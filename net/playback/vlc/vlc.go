package vlc

import (
	"fmt"
	"net/url"
	"os/exec"
	"simon/net/playback"
	"simon/pref"
	"strconv"
)

func run() {
	port := strconv.Itoa(pref.Get().VlcPort())
	vlcArgs := []string{"--extraintf", "http", "--http-port", port, "--http-password", "q"}
	cmd := exec.Command(Executable(), vlcArgs...)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

func Start() playback.Server {
	run()
	port := strconv.Itoa(pref.Get().VlcPort())
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port))
	server := New(addr)
	if err := server.Connect(); err != nil {
		panic(err)
	}
	return server
}
