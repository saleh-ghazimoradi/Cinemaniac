package middleware

import (
	"errors"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/domain"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/validator"
	"net/http"
	"strings"
)

func Authenticate(userRepo repository.UserRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = handlers.ContextSetUser(r, domain.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			helper.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if domain.ValidateTokenPlaintext(v, token); !v.Valid() {
			helper.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := userRepo.GetForToken(r.Context(), domain.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				helper.InvalidAuthenticationTokenResponse(w, r)
			default:
				helper.ServerErrorResponse(w, r, err)
			}
			return
		}
		r = handlers.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
