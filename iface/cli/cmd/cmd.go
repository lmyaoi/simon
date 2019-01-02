package cmd

import (
	"fmt"
	"net/url"
	"os"
	"simon/net/local"
	"simon/net/playback"
	"simon/net/playback/vlc"
	"simon/net/remote"
	"simon/conf"
	"strconv"
	"sync"
	"time"
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
	conf.Get().SetHostAddr(args[0])
	port, _ := strconv.Atoi(args[1])
	conf.Get().SetHostPort(port)
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
	port, _ := strconv.Atoi(args[0])
	conf.Get().SetVlcPort(port)
}

func _setvPath(args []string) {
	conf.Get().SetVlcPath(args[0])
}

func _setIval(args []string) {
	interval, _ := time.ParseDuration(args[0])
	conf.Get().SetInterval(interval)
}

func _status(args []string) {
	fmt.Printf("host address = http://%v:%v\n", conf.Get().HostAddr(), conf.Get().HostPort())
	fmt.Printf("vlc address  = http://localhost:%v\n", conf.Get().VlcPort())
	fmt.Printf("vlc path = \"%v\"\n", conf.Get().VlcPath())
	fmt.Printf("interval = %v\n", conf.Get().Interval())
}

func _save(args []string) {
	err := conf.Save()
	if err != nil {
		fmt.Println(err)
	}
}
