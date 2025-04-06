package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/middleware"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/repository"
	"net/http"
)

func movieRoutes(route *httprouter.Router, handler *handlers.MovieHandler, permission repository.PermissionRepository) {
	route.HandlerFunc(http.MethodPost, "/v1/movies", middleware.RequirePermission(permission, "movies:write", handler.CreateMovieHandler))
	route.HandlerFunc(http.MethodGet, "/v1/movies", middleware.RequirePermission(permission, "movies:read", handler.GetMoviesHandler))
	route.HandlerFunc(http.MethodGet, "/v1/movies/:id", middleware.RequirePermission(permission, "movies:read", handler.ShowMovieHandler))
	route.HandlerFunc(http.MethodPatch, "/v1/movies/:id", middleware.RequirePermission(permission, "movies:write", handler.UpdateMovieHandler))
	route.HandlerFunc(http.MethodDelete, "/v1/movies/:id", middleware.RequirePermission(permission, "movies:write", handler.DeleteMovieHandler))
}
