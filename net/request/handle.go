package request

import (
	"net/http"
	"vsync/net/playback"
)

type Handler struct {
	playback playback.Server
}

func Handle(playback playback.Server) *Handler {
	return &Handler{playback: playback}
}

func (h *Handler) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stat, err := h.playback.Status()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, _ = w.Write(stat.Marshal())
	}
}
