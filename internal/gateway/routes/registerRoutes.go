package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"net/http"
)

func RegisterRoutes() http.Handler {
	router := httprouter.New()

	healthHandler := handlers.NewHealthHandler()

	healthCheckRoutes(router, healthHandler)

	return router
}
