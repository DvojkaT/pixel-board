package handler

import (
	"board/internal/hub"
	"io/fs"
	"net/http"

	"github.com/centrifugal/centrifuge"
)

type Handler struct {
	hub    *hub.Hub
	static fs.FS
}

func NewHandler(hub *hub.Hub, static fs.FS) *Handler {
	return &Handler{
		hub:    hub,
		static: static,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	subFS, err := fs.Sub(h.static, "web")
	if err != nil {
		panic(err)
	}

	mux.Handle("/connection/websocket", centrifuge.NewWebsocketHandler(h.hub.Node(), centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}))
	mux.Handle("/", http.FileServer(http.FS(subFS)))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

}
