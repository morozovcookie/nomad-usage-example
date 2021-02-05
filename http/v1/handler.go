package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const (
	ContentTypeHeaderKey   = "Content-Type"
	ContentTypeHeaderValue = "application/json"
)

type Handler struct {
	router chi.Router

	logger *zap.Logger
}

type HandlerConfiguration interface {
	Logger() *zap.Logger
}

func NewHandler(cfg HandlerConfiguration) *Handler {
	return &Handler{
		router: chi.NewRouter(),

		logger: cfg.Logger(),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(ContentTypeHeaderKey, ContentTypeHeaderValue)

	h.router.ServeHTTP(w, r)
}
