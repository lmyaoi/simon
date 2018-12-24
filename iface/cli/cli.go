package cli

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"vsync/flags"
	"vsync/net/local"
	"vsync/net/playback"
	"vsync/net/playback/vlc"
	"vsync/net/remote"
)

type CLI int

func (CLI) Start() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`> `) // request a command
		cmd, _ := reader.ReadString('\n')
		handleCmd(cmd)
	}
}

func handleCmd(cmd string) {
	cmd = strings.TrimSpace(cmd)
	args := strings.Split(cmd, " ")
	if len(args) < 1 {
		return
	}
	switch args[0] {
	case "help", "list":
		help()
	case "host":
		host(args[1:])
	case "join":
		join(args[1:])
	case "pref":
		pref()
	case "exit", "quit":
		exit()
	}
}

func startVlc() {
	vlcArgs := []string{"--extraintf", "http", "--http-port", strconv.Itoa(*flags.VlcPort), "--http-password", "q"}
	cmd := exec.Command(flags.Vlc(), vlcArgs...)
	if err := cmd.Start(); err != nil {
		panic(1)
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

func host(args []string) {
	port, err := handleHostArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	server := startPlaybackServer()
	wg := &sync.WaitGroup{}
	host := local.NewHost(server, port)
	client := local.NewClient(playback.Dummy, host, wg)
	client.On()
	wg.Wait()
}

func handleHostArgs(args []string) (int, error) {
	if len(args) > 1 {
		return -1, errors.New("too many arguments, expected 1")
	} else if len(args) == 1 {
		return strconv.Atoi(args[0])
	}
	return *flags.HostPort, nil

}

func join(args []string) {
	addr, err := handleJoinArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	server := startPlaybackServer()
	wg := &sync.WaitGroup{}
	host := remote.NewHost(addr, vlc.Unmarshal)
	client := local.NewClient(server, host, wg)
	client.On()
	wg.Wait()
}

func handleJoinArgs(args []string) (*url.URL, error) {
	if len(args) > 1 {
		return nil, errors.New("too many arguments, expected 1")
	} else if len(args) == 1 {
		return url.Parse("http://" + args[0])
	}
	return url.Parse(fmt.Sprintf("http://%v:%v", *flags.HostUrl, *flags.HostPort))
}

func help() {
	fmt.Println(`help, list: List all the supported commands.`)
	fmt.Println(`join [host:port]: Join a connection.`)
	fmt.Println(`host [port]: Host a connection.`)
	fmt.Println(`pref: Modify the stored preferences.`)
	fmt.Println(`exit, quit: Exit the program.`)
}

func pref() {
	// todo: reading and writing preferences
	fmt.Println(`todo: reading and writing preferences`)
}

func exit() {
	os.Exit(0)
}