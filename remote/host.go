package remote

import (
	"net/http"
	"net/url"
	"vsync/playback"
)

type Host struct {
	client *http.Client
	url url.URL
	unmarshaler playback.StatusUnmarshaler
}

func (h *Host) Status() playback.Status {
	res, err := h.client.Get(h.url.String()+"/status")
	if err != nil {
		panic(err)
	}
	s := h.unmarshaler.Unmarshal(res.Body)
	defer res.Body.Close()
	return s
}