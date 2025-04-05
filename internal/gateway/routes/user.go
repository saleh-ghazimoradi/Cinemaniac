package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"net/http"
)

func userRoutes(route *httprouter.Router, handler *handlers.UserHandler) {
	route.HandlerFunc(http.MethodPost, "/v1/users", handler.RegisterUserHandler)
	route.HandlerFunc(http.MethodPut, "/v1/users/activated", handler.ActivateUserHandler)
	route.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", handler.CreateAuthenticationTokenHandler)
}
