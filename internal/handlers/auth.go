package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/angelchiav/go-ecommerce/internal/httpx"
	"github.com/angelchiav/go-ecommerce/internal/service"
	"github.com/angelchiav/go-ecommerce/internal/sqlc"
)

type Auth struct {
	auth *service.AuthService
	q    *sqlc.Queries
}

func NewAuth(auth *service.AuthService, q *sqlc.Queries) *Auth {
	return &Auth{auth: auth, q: q}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_json")
		return
	}

	id, err := h.auth.Register(r.Context(), req.Email, req.Password)

	if err == service.ErrEmailTaken {
		httpx.Error(w, http.StatusConflict, "email_taken")
		return
	}

	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_json")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err == service.ErrInvalidCreds {
		httpx.Error(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "server_error")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"access_token": token,
		"token_type":   "Bearer",
	})
}

func (h *Auth) Me(w http.ResponseWriter, r *http.Request) {
	userID := httpx.MustUserID(r)
	u, err := h.q.GetUserByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "server_error")
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"role":  u.Role,
	})
}
