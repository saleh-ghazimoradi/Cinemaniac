package handlers

import (
	"context"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(r *http.Request) *domain.User {
	user, ok := r.Context().Value(userContextKey).(*domain.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
