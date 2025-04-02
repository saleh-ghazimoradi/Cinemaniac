package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"net/http"
)

func userRoutes(route *httprouter.Router, handler *handlers.UserHandler) {
	route.HandlerFunc(http.MethodPost, "/v1/users", handler.RegisterUserHandler)
}
