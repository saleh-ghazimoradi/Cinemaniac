package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/middleware"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/service"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/transaction"
	"github.com/saleh-ghazimoradi/Cinemaniac/pkg/notification"
	"net/http"
)

func RegisterRoutes(db *sql.DB) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(helper.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.MethodNotAllowedResponse)

	movieRepository := repository.NewMovieRepository(db, db)
	userRepository := repository.NewUserRepository(db, db)
	tokenRepository := repository.NewTokenRepository(db, db)

	txService := transaction.NewTXService(db)
	SMTP, _ := notification.NewMailer(config.AppConfig.SMTP.Host, config.AppConfig.SMTP.Port, config.AppConfig.SMTP.UserName, config.AppConfig.SMTP.Password, config.AppConfig.SMTP.Sender)
	movieService := service.NewMovieService(movieRepository, txService)
	userService := service.NewUserService(userRepository, txService, SMTP, tokenRepository)

	healthHandler := handlers.NewHealthHandler()
	movieHandler := handlers.NewMovieHandler(movieService)
	userHandler := handlers.NewUserHandler(userService)

	healthCheckRoutes(router, healthHandler)
	movieRoutes(router, movieHandler)
	userRoutes(router, userHandler)

	return middleware.RecoverPanic(middleware.RateLimit(middleware.Authenticate(userRepository, router)))
}
