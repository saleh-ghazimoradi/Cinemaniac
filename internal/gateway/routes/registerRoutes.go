package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/middleware"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/transaction"
	"net/http"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(helper.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.MethodNotAllowedResponse)

	movieRepository := repository.NewMovieRepository(db, db)
	userRepository := repository.NewUserRepository(db, db)

	txService := transaction.NewTXService(db)
	movieService := service.NewMovieService(movieRepository, txService)
	userService := service.NewUserService(userRepository, txService)

	healthHandler := handlers.NewHealthHandler()
	movieHandler := handlers.NewMovieHandler(movieService)
	userHandler := handlers.NewUserHandler(userService)

	healthCheckRoutes(router, healthHandler)
	movieRoutes(router, movieHandler)
	userRoutes(router, userHandler)

	return middleware.RecoverPanic(middleware.RateLimit(router))
}
