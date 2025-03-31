package handlers

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"net/http"
)

type HealthHandler struct {
}

func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", config.AppConfig.Server.Env)
	fmt.Fprintf(w, "version: %s\n", config.AppConfig.Server.Version)
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
