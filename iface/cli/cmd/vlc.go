package cmd

import (
	"fmt"
	"net/url"
	"os/exec"
	"simon/flags"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"strconv"
)

func startVlc() {
	vlcArgs := []string{"--extraintf", "http", "--http-port", strconv.Itoa(*flags.VlcPort), "--http-password", "q"}
	cmd := exec.Command(flags.Vlc(), vlcArgs...)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

func startPlaybackServer() playback.Server {
	startVlc()
	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%v", *flags.VlcPort))
	server := vlc.New(addr)
	if err := server.Connect(); err != nil {
		panic(err)
	}
	return server
}
