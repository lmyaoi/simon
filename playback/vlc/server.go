package vlc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"vsync/flags"
	"vsync/log"
	"vsync/playback"
)

var (
	status = "/requests/status.json"
	//playlist = "/requests/playlist.json"
	play = status + "?command=pl_forceresume"
	pause = status + "?command=pl_forcepause"
	stop = status + "?command=pl_stop"
	jump = func (id int) string { return fmt.Sprintf(status + "?command=pl_play&id=%v", id) }
	seek = func (pos int64) string { return fmt.Sprintf(status + "?command=seek&val=%v", pos) }
)

type Server struct {
	addr *url.URL
	last *Status
	client *http.Client
	username, password string
	cmd *exec.Cmd
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
		state: playback.Stopped,
		pos: time.Unix(0, 0),
		created: time.Now(),
		id: -1,
	}
	return &Server{addr: addr, client: &http.Client{}, last: stat, username: "", password: *flags.VlcPwd, cmd: cmd}
}

func (vlc *Server) Connect() {
	if err := vlc.cmd.Start(); err != nil {
		panic(err)
	}

	req := vlc.newRequest("/")
	res, err := vlc.client.Do(req)
	for retries := 5; retries > 0 && err != nil; retries-- {
		time.Sleep(100 * time.Millisecond)
		log.Println("Retrying...")
		res, err = vlc.client.Do(req)
	}
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		panic("failed to connect to vlc interface")
	}
}

func (vlc *Server) SetState(s playback.State) {
	req := vlc.newRequest(getStatePath(s))
	res, _ := vlc.client.Do(req)
	defer vlc.Status()
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)

}

func getStatePath(s playback.State) string {
	switch s {
	case playback.Stopped: return stop
	case playback.Playing: return play
	case playback.Paused: return pause
	default: panic("unsupported playback state")
	}

}

func (vlc *Server) Start() {
	vlc.SetState(playback.Playing)
}

func (vlc *Server) Pause() {
	vlc.SetState(playback.Paused)
}

func (vlc *Server) Stop() {
	vlc.SetState(playback.Stopped)
}

func (vlc *Server) Sync(stat playback.Status) {
	s := verify(stat)
	//if vlc.last.id != s.id { vlc.jump(s.id) } todo: handle playlists

	if vlc.last.state != s.state {
		vlc.SetState(s.state)
		vlc.seek(playback.Now(s).Unix())
	} else if playback.WorthSeeking(vlc.last, s) {
		vlc.seek(playback.Now(s).Unix())
	}
}

func (vlc *Server) seek(s int64) {
	req := vlc.newRequest(seek(s))
	res, _ := vlc.client.Do(req)
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)
}

func (vlc *Server) jump(id int) {
	req := vlc.newRequest(jump(id))
	res, _ := vlc.client.Do(req)
	ioutil.ReadAll(res.Body)
	res.Body.Close()
}

func (vlc *Server) Status() playback.Status {
	if time.Now().Sub(vlc.last.created) < time.Second {
		return vlc.last
	}

	req := vlc.newRequest(status)
	res, _ := vlc.client.Do(req)
	defer res.Body.Close()

	s := NewStatus(res.Body)

	vlc.last = s
	return s
}

func (vlc *Server) Last() playback.Status {
	return vlc.last
}