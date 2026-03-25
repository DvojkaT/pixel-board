package handler

import (
	"board/internal/hub"
	"net/http"

	"github.com/centrifugal/centrifuge"
)

type Handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(h.hub.Node(), centrifuge.WebsocketConfig{}))
}
