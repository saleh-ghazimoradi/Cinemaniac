package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/middleware"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"net/http"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(helper.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.MethodNotAllowedResponse)

	movieRepository := repository.NewMovieRepository(db, db)

	movieService := service.NewMovieService(movieRepository)

	healthHandler := handlers.NewHealthHandler()
	movieHandler := handlers.NewMovieHandler(movieService)

	healthCheckRoutes(router, healthHandler)
	movieRoutes(router, movieHandler)

	return middleware.RecoverPanic(router)
}
