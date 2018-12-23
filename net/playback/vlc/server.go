package vlc

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"vsync/net/httputil"
	"vsync/log"
	"vsync/net/playback"
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
		&jsonStatus{
			State:   playback.Stopped,
			Pos:     time.Unix(0, 0),
			Created: time.Now(),
			Id:      -1,
		},
	}
	return &Server{addr: addr, client: &http.Client{}, last: stat, username: "", password: "q", cmd: cmd}
}

func (vlc *Server) Connect() error {
	req := newRequest(vlc, status)
	res, err := httputil.Retry(vlc.client, req, 10)
	if err != nil {
		return err
	}
	vlc.last = NewStatus(res.Body)
	return nil
}


func commandify(state playback.State) string {
	switch state {
	case playback.Playing: return play
	case playback.Paused: return pause
	case playback.Stopped: return stop
	default:
		panic("Unsupported state")
	}
}

func (vlc *Server) SetState(s playback.State) error {
	req := newRequest(vlc, commandify(s))
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

	if vlc.last.State() != s.State() {
		if err := vlc.SetState(s.State()); err != nil {
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
	if !vlc.staleLast() {
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

func (vlc *Server) staleLast() bool {
	return time.Now().Sub(vlc.last.Created()) > time.Second
}

func (vlc *Server) Last() playback.Status {
	return vlc.last
}
