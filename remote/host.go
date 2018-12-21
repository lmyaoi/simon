package remote

import (
	"net/http"
	"net/url"
	"time"
	"vsync/playback"
)

type Host struct {
	client    *http.Client
	url       *url.URL
	unmarshal playback.StatusUnmarshaler
}

func NewHost(url *url.URL, unmarshal playback.StatusUnmarshaler) *Host {
	return &Host{&http.Client{Timeout: 10 * time.Second}, url, unmarshal}
}

func (h *Host) Status() (playback.Status, error) {
	res, err := h.client.Get(h.url.String() + "/status")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	s, err := h.unmarshal(res.Body)
	if err != nil {
		return nil, err
	}
	return s, nil
}
