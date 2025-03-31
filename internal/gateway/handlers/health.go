package handlers

import (
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"net/http"
)

type HealthHandler struct {
}

func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": config.AppConfig.Server.Env,
		"version":     config.AppConfig.Server.Version,
	}

	if err := helper.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		slg.Logger.Error(err.Error())
		http.Error(w, "The server encountered an internal error", http.StatusInternalServerError)
	}
}
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
