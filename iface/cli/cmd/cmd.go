package cmd

import (
	"fmt"
	"net/url"
	"os"
	"simon/flags"
	"simon/net/local"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"simon/net/remote"
	"strconv"
	"sync"
)

var (
	help    = newCmd(_help, check(noArgs))
	pref    = newCmd(_pref, check(noArgs))
	exit    = newCmd(_exit, check(noArgs))
	join    = newCmd(_join, defaultArgs(*flags.HostUrl, strconv.Itoa(*flags.HostPort)), check(validUrl))
	host    = newCmd(_host, defaultArgs(strconv.Itoa(*flags.HostPort)), check(validPort))
	vlcPort = newCmd(_vlcPort, defaultArgs(strconv.Itoa(*flags.VlcPort)), check(validPort))
	status = newCmd(_status, check(noArgs))
)

var Cmds = map[string]Runner{
	"help":     help,
	"list":     help,
	"host":     host,
	"join":     join,
	"exit":     exit,
	"quit":     exit,
	"status":	status,
	"vlc-port": vlcPort,
	"pref":     pref,
}

type Runner interface {
	Run([]string)
}

func _help([]string) {
	fmt.Println("help, list: List all the supported commands")
	fmt.Println("join [host port]: Join a connection")
	fmt.Println("host [port]: Host a connection")
	fmt.Println("vlc-port <port>: Set the vlc port to host the playback server at")
	fmt.Println("status: Show the current environment variables")
	fmt.Println("pref: Modify the stored preferences")
	fmt.Println("exit, quit: Exit the program")
}

func _pref([]string) {
	// todo: reading and writing preferences
	fmt.Println(`todo: reading and writing preferences`)
}

func _exit([]string) {
	os.Exit(0)
}

func _join(args []string) {
	addr, _ := url.Parse(fmt.Sprintf("http://%v:%v", args[0], args[1]))
	server := startPlaybackServer()
	wg := &sync.WaitGroup{}
	host := remote.NewHost(addr, vlc.Unmarshal)
	client := local.NewClient(server, host, wg)
	client.On()
	wg.Wait()
}

func _host(args []string) {
	port, _ := strconv.Atoi(args[0])
	server := startPlaybackServer()
	wg := &sync.WaitGroup{}
	host := local.NewHost(server, port)
	client := local.NewClient(playback.Dummy, host, wg)
	client.On()
	wg.Wait()
}

func _vlcPort(args []string) {
	*flags.VlcPort, _ = strconv.Atoi(args[0])
}

func _status(args []string) {
	fmt.Printf("host address = http://%v:%v\n", *flags.HostUrl, *flags.HostPort)
	fmt.Printf("vlc address  = http://localhost:%v\n", *flags.VlcPort)
	fmt.Printf("vlc path = \"%v\"\n", *flags.VlcPath)
	fmt.Printf("interval = %v\n", *flags.Interval)
}
