package local

import (
	"fmt"
	"net/http"
	"vsync/flags"
	"vsync/playback"
)

type Signal int

const (
	On Signal = iota
	Kill
)

type Host struct {
	playback playback.Server
	server   http.Server
	control  <-chan Signal
}

func (h *Host) Status() (playback.Status, error) {
	return h.playback.Last(), nil // when client asks tell it that its server is up to date since they are both local
}

func (h *Host) start() {
	panic(h.server.ListenAndServe())
}

func NewHost(playback playback.Server) *Host {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", *flags.HostPort),
		Handler: mux,
	}
	mux.HandleFunc("/status", serveStatus(playback))
	control := make(chan Signal)
	h := &Host{playback, server, control}
	go h.start()
	return h
}

func serveStatus(playback playback.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stat, err := playback.Status()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, _ = w.Write(stat.Marshal())
	}
}
