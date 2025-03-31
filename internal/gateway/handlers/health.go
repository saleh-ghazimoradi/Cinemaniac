package handlers

import (
	"github.com/saleh-ghazimoradi/Cinemaniac/config"
	"github.com/saleh-ghazimoradi/Cinemaniac/internal/helper"
	"net/http"
)

type HealthHandler struct {
}

func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := helper.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": config.AppConfig.Server.Env,
			"version":     config.AppConfig.Server.Version,
		},
	}

	if err := helper.WriteJSON(w, http.StatusOK, env, nil); err != nil {
		helper.ServerErrorResponse(w, r, err)
	}
}
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
