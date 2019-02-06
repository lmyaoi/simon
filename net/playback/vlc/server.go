package vlc

import (
	"fmt"
	"net/http"
	"net/url"
	"simon/conf"
	"simon/log"
	"simon/net/httputil"
	"simon/net/playback"
	"simon/ticker"
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
	polled             *playback.Buffer
	offered            *playback.Buffer
	client             *http.Client
	username, password string
	control            chan<- ticker.Signal
	signal             <-chan time.Time
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
	signal, control := ticker.New(conf.Get().Interval.Duration)
	s := &Server{
		addr:     addr,
		client:   &http.Client{},
		polled:   &playback.Buffer{},
		offered:  &playback.Buffer{},
		username: "",
		password: "q",
		control:  control,
		signal:   signal,
	}
	go s.loop()
	return s
}

func (server *Server) Connect() error {
	req := newRequest(server, status)
	res, err := httputil.Retry(server.client, req, 10)
	if err != nil {
		return err
	}
	defer httputil.Discard(res, err)
	server.polled.Push(NewStatus(res.Body))
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

func (server *Server) setState(s playback.State) error {
	req := newRequest(server, commandify(s))
	res, err := server.client.Do(req)
	if err != nil {
		return err
	}
	defer httputil.Discard(res, err)

	stat := NewStatus(res.Body)
	go server.polled.Push(stat) //update last polled status
	stat = stat.Copy()
	stat.SetState(s)
	go server.offered.Push(stat) //update last offered status
	return nil
}

func (server *Server) Start() error {
	return server.setState(playback.Playing)
}

func (server *Server) Pause() error {
	return server.setState(playback.Paused)
}

func (server *Server) Stop() error {
	return server.setState(playback.Stopped)
}

func (server *Server) forceSync(s playback.Status) error {
	if s == nil {
		return nil
	}
	status := Verify(s)
	// todo: handle playlists
	latest := server.Last()
	if latest == nil || latest.State() != status.State() {
		if err := server.setState(status.State()); err != nil {
			return err
		}
		server.seek(playback.Now(status).Unix())
	} else if playback.WorthSeeking(latest, status) {
		server.seek(playback.Now(status).Unix())
	}
	return nil
}

func (server *Server) Sync(s playback.Status) error {
	go server.offered.Push(s)
	return nil
}

func (server *Server) loop() {
	for {
		select {
		case <-server.signal:
			if err := server.forceSync(server.offered.Peek()); err != nil {
				log.Println(err)
			}
		}
	}
}

func (server *Server) seek(s int64) {
	req := newRequest(server, seek(s))

	res, err := server.client.Do(req)
	if err != nil {
		return
	}
	defer httputil.Discard(res, err)

	stat := NewStatus(res.Body)
	go server.polled.Push(stat) //update last polled status
	stat = stat.Copy()
	stat.SetPos(s)
	go server.offered.Push(stat) //update last offered status
}

func (server *Server) jump(id int) {
	req := newRequest(server, jump(id))

	res, err := server.client.Do(req)
	if err != nil {
		return
	}
	defer httputil.Discard(res, err)

	stat := NewStatus(res.Body)
	go server.polled.Push(stat) //update last polled status
	stat = stat.Copy()
	stat.SetId(id)
	go server.offered.Push(stat) //update last offered status
}

func (server *Server) Status() (playback.Status, error) {
	polledStatus := server.polled.Peek()
	if !playback.Stale(polledStatus) {
		return polledStatus, nil
	}

	req := newRequest(server, status)
	res, err := server.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httputil.Discard(res, err)

	status := NewStatus(res.Body)
	server.polled.Push(status)
	return status, nil
}

func (server *Server) Last() playback.Status {
	polled := server.polled.Peek()
	offered := server.offered.Peek()
	if offered == nil || polled == nil || offered.Created().After(polled.Created()) {
		return polled
	}
	return offered
}

func (server *Server) On() {
	server.control <- ticker.On
}

func (server *Server) Off() {
	server.control <- ticker.Off
}

func (server *Server) Kill() {
	server.control <- ticker.Kill
}
