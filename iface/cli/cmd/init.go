package cmd

import (
	"fmt"
	"net/url"
	"os"
	"simon/net/local"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"simon/net/remote"
	"simon/pref"
	"strconv"
	"sync"
)

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
	server := vlc.Start()
	wg := &sync.WaitGroup{}
	host := remote.NewHost(addr, vlc.Unmarshal)
	client := local.NewClient(server, host, wg)
	client.On()
	wg.Wait()
}

func _host(args []string) {
	port, _ := strconv.Atoi(args[0])
	server := vlc.Start()
	wg := &sync.WaitGroup{}
	host := local.NewHost(server, port)
	client := local.NewClient(playback.Dummy, host, wg)
	client.On()
	wg.Wait()
}

func _vlcPort(args []string) {
	port, _ := strconv.Atoi(args[0])
	pref.Get().SetVlcPort(port)
}

func _status(args []string) {
	fmt.Printf("host address = http://%v:%v\n", pref.Get().HostUrl(), pref.Get().HostPort())
	fmt.Printf("vlc address  = http://localhost:%v\n", pref.Get().VlcPort())
	fmt.Printf("vlc path = \"%v\"\n", pref.Get().VlcPath())
	fmt.Printf("interval = %v\n", pref.Get().Interval())
}
