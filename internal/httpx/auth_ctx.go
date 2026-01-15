package httpx

import (
	"context"
	"net/http"
)

type CtxKey string

const userIDKey CtxKey = "user_id"
const userRoleKey CtxKey = "user_role"

func WithAuth(ctx context.Context, userID int64, role string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, userRoleKey, role)
	return ctx
}

func UserID(r *http.Request) (int64, bool) {
	v := r.Context().Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}

func MustUserID(r *http.Request) int64 {
	id, ok := UserID(r)
	if !ok {
		panic("missing user id in context")
	}
	return id
}
