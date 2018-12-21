package vlc

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"vsync/httputil"
	"vsync/log"
	"vsync/playback"
)

var (
	status = "/requests/status.json"
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

func newRequest(vlc *Server, path string) *http.Request {
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

	req := newRequest(vlc, status)
	res, err := vlc.retry(req, 10)
	if err != nil {
		return err
	}
	vlc.last = NewStatus(res.Body)
	return nil
}

func (vlc *Server) retry(req *http.Request, retries int) (res *http.Response, err error) {
	for i := 0; i < 1 + retries; i++ {
		res, err = vlc.client.Do(req)
		if err == nil { return }
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (vlc *Server) SetState(s playback.State) error {
	req := newRequest(vlc, s.String())
	res, err := vlc.client.Do(req)
	if err != nil {
		return err
	}
	defer httputil.Discard(res, err)
	defer func() {
		if _, err := vlc.Status(); err != nil {
			log.Println(err)
		}
	}()
	return nil
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
	// todo: handle playlists

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
	req := newRequest(vlc, seek(s))
	httputil.Discard(vlc.client.Do(req))
}

func (vlc *Server) jump(id int) {
	req := newRequest(vlc, jump(id))
	httputil.Discard(vlc.client.Do(req))
}

func (vlc *Server) Status() (playback.Status, error) {
	if time.Now().Sub(vlc.last.created) < time.Second {
		return vlc.last, nil
	}

	req := newRequest(vlc, status)
	res, err := vlc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httputil.Discard(res, err)

	vlc.last = NewStatus(res.Body)
	return vlc.last, nil
}

func (vlc *Server) Last() playback.Status {
	return vlc.last
}
