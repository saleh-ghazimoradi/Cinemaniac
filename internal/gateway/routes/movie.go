package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"net/http"
)

func movieRoutes(route *httprouter.Router, handler *handlers.MovieHandler) {
	route.HandlerFunc(http.MethodPost, "/v1/movies", handler.CreateMovieHandler)
	route.HandlerFunc(http.MethodGet, "/v1/movies/:id", handler.ShowMovieHandler)
	route.HandlerFunc(http.MethodPut, "/v1/movies/:id", handler.UpdateMovieHandler)
}
