package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"simon/flags"
	"simon/net"
	"simon/net/local"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"simon/net/remote"
	"strconv"
	"sync"
)

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}

	addr, _ := url.Parse(fmt.Sprintf("http://localhost:%v", *flags.VlcPort))

	vlcArgs := []string{"--extraintf", "http", "--http-port", strconv.Itoa(*flags.VlcPort), "--http-password", "q"}
	cmd := exec.Command(flags.Vlc(), vlcArgs...)
	server := vlc.New(addr, cmd)
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	if err := server.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalln(err)
		}
	}()

	var h net.Host
	var client *local.Client
	if *flags.Host {
		h = local.NewHost(server)
		client = local.NewClient(playback.Dummy, h, wg)
	} else {
		hostAddr, _ := url.Parse(fmt.Sprintf("http://%v:%v", *flags.HostUrl, *flags.HostPort))
		h = remote.NewHost(hostAddr, vlc.Unmarshal)
		client = local.NewClient(server, h, wg)
	}
	client.On()
	wg.Wait()
}
