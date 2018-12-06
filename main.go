package main

import (
	"flag"
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"sync"
	"vsync/flags"
	"vsync/host"
	"vsync/local"
	"vsync/log"
	"vsync/playback"
	"vsync/playback/vlc"
	"vsync/remote"
)

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}

	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%v", *flags.VlcPort))

	vlcArgs := []string{"--extraintf", "http", "--http-port", strconv.Itoa(*flags.VlcPort), "--http-password", "q"}
	cmd := exec.Command(flags.Vlc(), vlcArgs...)
	server := vlc.New(addr, cmd)

	if err := server.Connect(); err != nil {
		log.Println(err)
	}
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Println(err)
		}
	}()

	var h host.Host
	var client *local.Client
	if *flags.Host {
		h = local.NewHost(server)
		client = local.NewClient(playback.D, h, wg)
	} else {
		hostAddr, _ := url.Parse(fmt.Sprintf("http://%v:%v", *flags.HostUrl, *flags.HostPort))
		h = remote.NewHost(hostAddr, vlc.Unmarshal)
		client = local.NewClient(server, h, wg)
	}
	client.On()
	wg.Wait()
}
