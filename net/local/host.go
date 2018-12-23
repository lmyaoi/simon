package local

import (
	"fmt"
	"net/http"
	"vsync/flags"
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
	addr := fmt.Sprintf("%v:%v", *flags.HostUrl, *flags.HostPort)
	return &http.Server{Addr: addr, Handler: mux}
}