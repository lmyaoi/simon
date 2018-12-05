package vlc

import (
	"fmt"
	"net/url"
	"os/exec"
	"testing"
	"time"
	"vsync/playback"
)

func setup() (*Server, *exec.Cmd) {
	addr, _ := url.Parse("http://:9090")
	server := New(addr)
	cmd := exec.Command("/Applications/VLC.app/Contents/MacOS/VLC", "/Users/rensux/Downloads/test.mkv", "--extraintf", "http", "--http-port", "9090", "--http-password", "q")
	return server, cmd
}

func TestServer_Connect(t *testing.T) {
	server, cmd := setup()
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	defer cmd.Process.Kill()
	server.Connect()
}

func TestServer_Start(t *testing.T) {
	server, cmd := setup()
	cmd.Start()
	defer cmd.Process.Kill()
	server.Connect()
	time.Sleep(3 * time.Second)
	server.SetState(playback.Paused)
	time.Sleep(2 * time.Second)
	server.SetState(playback.Playing)
	time.Sleep(2 * time.Second)
}

func TestServer_Status(t *testing.T) {
	server, cmd := setup()
	cmd.Start()
	defer cmd.Process.Kill()

	server.Connect()
	time.Sleep(3 * time.Second)
	stat := verify(server.Status())
	fmt.Printf("id: %v, state: %v, pos: %v, created: %v\n", stat.id, stat.State(), stat.pos.Unix(), stat.created.Unix())
}

func TestServer_Sync(t *testing.T) {
	server, cmd := setup()
	cmd.Start()
	defer cmd.Process.Kill()

	server.Connect()
	time.Sleep(3 * time.Second)
	server.Pause()
	stat := *verify(server.Status())
	fmt.Printf("id: %v, state: %v, pos: %v, created: %v\n", stat.id, stat.State(), stat.pos.Unix(), stat.created.Unix())
	server.Start()

	time.Sleep(5 * time.Second)
	stat.state = playback.Playing
	server.Sync(&stat)
	stat = *verify(server.Status())
	fmt.Printf("id: %v, state: %v, pos: %v, created: %v\n", stat.id, stat.State(), stat.pos.Unix(), stat.created.Unix())
	time.Sleep(2 * time.Second)
}
