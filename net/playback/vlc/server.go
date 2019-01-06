package vlc

import (
	"fmt"
	"net/http"
	"net/url"
	"simon/log"
	"simon/net/httputil"
	"simon/net/playback"
	"time"
)

var (
	status = "/requests/status.json"
	play   = status + "?command=pl_forceresume"
	pause  = status + "?command=pl_forcepause"
	stop   = status + "?command=pl_stop"
	jump   = func(id int) string { return fmt.Sprintf(status+"?command=pl_play&id=%v", id) }
	seek   = func(pos int64) string { return fmt.Sprintf(status+"?command=seek&val=%v", pos) }
)

type Server struct {
	addr               *url.URL
	last               *Status
	client             *http.Client
	username, password string
}

func newRequest(vlc *Server, path string) *http.Request {
	req, err := http.NewRequest("GET", vlc.addr.String()+path, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(vlc.username, vlc.password)
	return req
}

func New(addr *url.URL) *Server {
	stat := DefaultStatus
	return &Server{addr: addr, client: &http.Client{}, last: &stat, username: "", password: "q"}
}

func (server *Server) Connect() error {
	req := newRequest(server, status)
	res, err := httputil.Retry(server.client, req, 10)
	if err != nil {
		return err
	}
	server.last = NewStatus(res.Body)
	return nil
}

func commandify(state playback.State) string {
	switch state {
	case playback.Playing:
		return play
	case playback.Paused:
		return pause
	case playback.Stopped:
		return stop
	default:
		panic("Unsupported state")
	}
}

func (server *Server) SetState(s playback.State) error {
	req := newRequest(server, commandify(s))
	res, err := server.client.Do(req)
	if err != nil {
		return err
	}
	defer httputil.Discard(res, err)
	defer func() {
		if _, err := server.Status(); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func (server *Server) Start() error {
	return server.SetState(playback.Playing)
}

func (server *Server) Pause() error {
	return server.SetState(playback.Paused)
}

func (server *Server) Stop() error {
	return server.SetState(playback.Stopped)
}

func (server *Server) Sync(stat playback.Status) error {
	s := Verify(stat)
	// todo: handle playlists

	if server.last.State() != s.State() {
		if err := server.SetState(s.State()); err != nil {
			return err
		}
		server.seek(playback.Now(s).Unix())
	} else if playback.WorthSeeking(server.last, s) {
		server.seek(playback.Now(s).Unix())
	}
	return nil
}

func (server *Server) seek(s int64) {
	req := newRequest(server, seek(s))
	httputil.Discard(server.client.Do(req))
}

func (server *Server) jump(id int) {
	req := newRequest(server, jump(id))
	httputil.Discard(server.client.Do(req))
}

func (server *Server) Status() (playback.Status, error) {
	if !server.staleLast() {
		return server.last, nil
	}

	req := newRequest(server, status)
	res, err := server.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httputil.Discard(res, err)

	server.last = NewStatus(res.Body)
	return server.last, nil
}

func (server *Server) staleLast() bool {
	return time.Now().Sub(server.last.Created()) > time.Second
}

func (server *Server) Last() playback.Status {
	return server.last
}
