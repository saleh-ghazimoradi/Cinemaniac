package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/gateway/handlers"
	"net/http"
)

func healthCheckRoutes(route *httprouter.Router, handler *handlers.HealthHandler) {
	route.HandlerFunc(http.MethodGet, "/v1/healthcheck", handler.HealthCheckHandler)
}
