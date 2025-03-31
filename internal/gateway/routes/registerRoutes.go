package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"net/http"
)

func RegisterRoutes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(helper.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.MethodNotAllowedResponse)

	healthHandler := handlers.NewHealthHandler()
	movieHandler := handlers.NewMovieHandler()

	healthCheckRoutes(router, healthHandler)
	movieRoutes(router, movieHandler)

	return router
}
