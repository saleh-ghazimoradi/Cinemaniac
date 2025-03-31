package handlers

import "net/http"

type HealthHandler struct {
}

func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"alive": true}`))
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
