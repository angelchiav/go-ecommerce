package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/angelchiav/go-ecommerce/internal/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCreds = errors.New("invalid_credentials")
	ErrEmailTaken   = errors.New("email_taken")
)

type AuthService struct {
	q      *sqlc.Queries
	secret []byte
}

func NewAuthService(q *sqlc.Queries, jwtSecret string) *AuthService {
	return &AuthService{q: q, secret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int64, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || len(password) < 8 {
		return 0, errors.New("email_required_password_min_8")
	}

	_, err := s.q.GetUserByEmail(ctx, email)
	if err == nil {
		return 0, ErrEmailTaken
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	u, err := s.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return 0, err
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.q.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return "", ErrInvalidCreds
	}
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCreds
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(u.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, struct {
		Role string `json:"role"`
		jwt.RegisteredClaims
	}{
		Role:             u.Role,
		RegisteredClaims: claims,
	})

	return t.SignedString(s.secret)
}
