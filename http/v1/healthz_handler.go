package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	HealthzHandlerPathPrefix = "/v1/healthz"

	HealthzStatusUp = "UP"
)

type HealthzHandlerConfiguration interface {
	HandlerConfiguration
}

type HealthzHandler struct {
	*Handler

	upTime time.Time
}

func NewHealthzHandler(cfg HealthzHandlerConfiguration) *HealthzHandler {
	h := &HealthzHandler{
		Handler: NewHandler(cfg),

		upTime: time.Now(),
	}

	h.router.Get(HealthzHandlerPathPrefix, h.handleHealthz)

	return h
}

type HealthzResponse struct {
	Status string `json:"status"`
	UpTime string `json:"up_time"`
}

func encodeHealthzResponse(_ context.Context, w http.ResponseWriter, status int, resp *HealthzResponse) error {
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(resp)
}

func (h *HealthzHandler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := &HealthzResponse{
		Status: HealthzStatusUp,
		UpTime: time.Since(h.upTime).Round(time.Millisecond).String(),
	}

	if err := encodeHealthzResponse(r.Context(), w, http.StatusOK, resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
