package remote

import (
	"net/http"
	"net/url"
	"vsync/playback"
)

type Host struct {
	client *http.Client
	url *url.URL
	unmarshal playback.StatusUnmarshaler
}

func NewHost(url *url.URL, unmarshal playback.StatusUnmarshaler) *Host {
	return &Host{&http.Client{}, url, unmarshal}
}

func (h *Host) Status() playback.Status {
	res, err := h.client.Get(h.url.String()+"/status")
	if err != nil {
		panic(err)
	}
	s := h.unmarshal(res.Body)

	defer res.Body.Close()
	return s
}