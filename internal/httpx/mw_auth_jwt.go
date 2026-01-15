package httpx

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func AuthJWT(secret string) func(next http.HandlerFunc) http.HandlerFunc {
	sec := []byte(secret)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				Error(w, http.StatusUnauthorized, "missing_bearer_token")
				return
			}
			tokenStr := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))

			claims := &Claims{}
			tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return sec, nil
			})
			if err != nil || !tok.Valid {
				Error(w, http.StatusUnauthorized, "invalid_token")
				return
			}

			uid, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil || uid <= 0 {
				Error(w, http.StatusUnauthorized, "invalid_token_subject")
				return
			}

			ctx := WithAuth(r.Context(), uid, claims.Role)
			next(w, r.WithContext(ctx))
		}
	}
}
