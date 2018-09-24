package host

import (
	"fmt"
	"net/http"
	"vsync/flags"
	"vsync/vlc"
	"vsync/log"
)

type Host struct {
	s http.Server
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("update request received from %v\n", r.RemoteAddr)
	fmt.Fprintf(w, "%v\n", vlc.Latest.String())
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("connection request received from %v\n", r.RemoteAddr)
	fmt.Fprintf(w, "%v\n", "welcome to the darkside love")
}

func New() *Host {
	url := fmt.Sprintf(":%v", *flags.Port)
	s := http.Server{Addr: url}
	h := &Host{s}
	return h
}

func (h *Host) Listen() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/connect", connectionHandler)
	h.s.ListenAndServe()
}