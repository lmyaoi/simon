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
	"vsync/playback"
	"vsync/playback/vlc"
	"vsync/remote"
)

func main() {
	flag.Parse()
	wg := &sync.WaitGroup{}

	addr, _ := url.Parse("http://:" + vlcPort())
	cmd := exec.Command("/Applications/VLC.app/Contents/MacOS/VLC", "--extraintf", "http", "--http-port", vlcPort(), "--http-password", *flags.VlcPwd)
	server := vlc.New(addr, cmd)

	server.Connect()
	defer cmd.Process.Kill()

	var h host.Host
	var client *local.Client
	if *flags.Host {
		h = local.NewHost(server)
		client = local.NewClient(playback.D, h, wg)
	} else {
		hostAddr, _ := url.Parse(fmt.Sprintf("http://:%v", *flags.HostPort))
		h = remote.NewHost(hostAddr, vlc.Unmarshal)
		client = local.NewClient(server, h, wg)
	}
	client.On()
	wg.Wait()
}

func vlcPort() string {
	if *flags.Host {
		return strconv.FormatInt(int64(*flags.VlcPort), 10)
	}
	return strconv.FormatInt(int64(*flags.VlcPort + 1), 10)
}