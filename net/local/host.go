package local

import (
	"fmt"
	"net/http"
	"vsync/net/playback"
	"vsync/net/request"
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

func NewHost(playback playback.Server, port int) *Host {
	server := newServer(request.Handle(playback), port)
	control := make(chan Signal)
	h := &Host{playback, server, control}
	go h.start()
	return h
}

func newServer(handler *request.Handler, port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", handler.Status())
	addr := fmt.Sprintf("localhost:%v", port)
	return &http.Server{Addr: addr, Handler: mux}
}
