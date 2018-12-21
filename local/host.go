package local

import (
	"net/http"
	"vsync/consts"
	"vsync/local/request"
	"vsync/playback"
)

//go:generate stringer -type=Signal
type Signal int

const (
	On Signal = iota
	Kill
)

type Host struct {
	playback playback.Server
	server   *http.Server
	control  <-chan Signal
}

func (h *Host) Status() (playback.Status, error) {
	return h.playback.Last(), nil
}

func (h *Host) start() {
	panic(h.server.ListenAndServe())
}

func NewHost(playback playback.Server) *Host {
	server := newServer(request.Handle(playback))
	control := make(chan Signal)
	h := &Host{playback, server, control}
	go h.start()
	return h
}

func newServer(handler *request.Handler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", handler.Status())
	return &http.Server{Addr: consts.HostAddr, Handler: mux}
}