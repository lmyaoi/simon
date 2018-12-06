package vlc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"vsync/log"
	"vsync/playback"
)

var (
	status = "/requests/status.json"
	//playlist = "/requests/playlist.json"
	play  = status + "?command=pl_forceresume"
	pause = status + "?command=pl_forcepause"
	stop  = status + "?command=pl_stop"
	jump  = func(id int) string { return fmt.Sprintf(status+"?command=pl_play&id=%v", id) }
	seek  = func(pos int64) string { return fmt.Sprintf(status+"?command=seek&val=%v", pos) }
)

type Server struct {
	addr               *url.URL
	last               *Status
	client             *http.Client
	username, password string
	cmd                *exec.Cmd
}

func (vlc *Server) newRequest(path string) *http.Request {
	req, err := http.NewRequest("GET", vlc.addr.String()+path, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(vlc.username, vlc.password)
	return req
}

func New(addr *url.URL, cmd *exec.Cmd) *Server {
	stat := &Status{
		state:   playback.Stopped,
		pos:     time.Unix(0, 0),
		created: time.Now(),
		id:      -1,
	}
	return &Server{addr: addr, client: &http.Client{}, last: stat, username: "", password: "q", cmd: cmd}
}

func (vlc *Server) Connect() error {
	if err := vlc.cmd.Start(); err != nil {
		return err
	}

	req := vlc.newRequest(status)
	res, err := vlc.client.Do(req)
	for i := 0; i < 10 && err != nil; i++ {
		log.Print("Failed to connect to playback server. Retrying...")
		time.Sleep(100 * time.Millisecond)
		res, err = vlc.client.Do(req)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	vlc.last = NewStatus(res.Body)
	return nil
}

func (vlc *Server) SetState(s playback.State) error {
	req := vlc.newRequest(getStatePath(s))
	res, err := vlc.client.Do(req)
	if err != nil {
		return err
	}
	defer vlc.Status()
	defer res.Body.Close()
	defer ioutil.ReadAll(res.Body)
	return nil
}

func getStatePath(s playback.State) string {
	switch s {
	case playback.Stopped:
		return stop
	case playback.Playing:
		return play
	case playback.Paused:
		return pause
	default:
		panic("Unsupported playback state.")
	}
}

func (vlc *Server) Start() error {
	return vlc.SetState(playback.Playing)
}

func (vlc *Server) Pause() error {
	return vlc.SetState(playback.Paused)
}

func (vlc *Server) Stop() error {
	return vlc.SetState(playback.Stopped)
}

func (vlc *Server) Sync(stat playback.Status) error {
	s := verify(stat)
	//if vlc.last.id != s.id { vlc.jump(s.id) } todo: handle playlists

	if vlc.last.state != s.state {
		if err := vlc.SetState(s.state); err != nil {
			return err
		}
		vlc.seek(playback.Now(s).Unix())
	} else if playback.WorthSeeking(vlc.last, s) {
		vlc.seek(playback.Now(s).Unix())
	}
	return nil
}

func (vlc *Server) seek(s int64) {
	req := vlc.newRequest(seek(s))
	res, err := vlc.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	defer ioutil.ReadAll(res.Body)
}

func (vlc *Server) jump(id int) {
	req := vlc.newRequest(jump(id))
	res, err := vlc.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	defer ioutil.ReadAll(res.Body)
}

func (vlc *Server) Status() (playback.Status, error) {
	if time.Now().Sub(vlc.last.created) < time.Second {
		return vlc.last, nil
	}

	req := vlc.newRequest(status)
	res, err := vlc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	vlc.last = NewStatus(res.Body)
	return vlc.last, nil
}

func (vlc *Server) Last() playback.Status {
	return vlc.last
}