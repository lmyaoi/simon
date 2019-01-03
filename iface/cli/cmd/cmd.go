package cmd

import (
	"fmt"
	"net/url"
	"os"
	"simon/conf"
	"simon/jsonutil"
	"simon/net/local"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"simon/net/remote"
	"strconv"
	"sync"
)

type Runner interface {
	Run([]string)
}

func _help([]string) {
	fmt.Println("help, list: List all the supported commands")
	fmt.Println("join [host port]: Join a connection")
	fmt.Println("sethost <host port>: Set the default host address to join")
	fmt.Println("setival <interval>: Set the interval at which to poll")
	fmt.Println("host [port]: Host a connection")
	fmt.Println("setvport <port>: Set the vlc port to host the playback server at")
	fmt.Println("setvpath <path>: Set the path of the vlc executable")
	fmt.Println("status: Show the current environment variables")
	fmt.Println("save: Save the current in memory config")
	fmt.Println("exit, quit: Exit the program")
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

func _setHost(args []string) {
	conf.Get().HostAddr = args[0]
	conf.Get().HostPort, _ = strconv.Atoi(args[1])
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

func _setvPort(args []string) {
	conf.Get().VlcPort, _ = strconv.Atoi(args[0])
}

func _setIval(args []string) {
	conf.Get().Interval, _ = jsonutil.ParseDuration(args[0])
}

func _status(args []string) {
	fmt.Println(conf.Get())
}

func _save(args []string) {
	err := conf.Save()
	if err != nil {
		fmt.Println(err)
	}
}
