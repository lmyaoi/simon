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
	server http.Server
	control <-chan Signal
}

func (h *Host) Status() playback.Status {
	return h.playback.Last() // when client asks tell it that its server is up to date since they are both local
}

func NewHost(playback playback.Server) *Host {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", serveStatus(playback))
	server := http.Server{
		Addr:              fmt.Sprintf(":%v", *flags.HostPort),
		Handler:           mux,
	}
	control := make(chan Signal)
	h := &Host{playback, server, control}
	go h.server.ListenAndServe()
	return h
}

func serveStatus(server playback.Server) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		stat := server.Status().Marshal()
		if _, err := w.Write(stat); err != nil {
			panic(err)
		}
	}
}