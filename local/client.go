package local

import (
	"vsync/host"
	"vsync/playback"
)

type Client struct {
	server playback.Server
	host   host.Host
}

func NewClient(s playback.Server, h host.Host) *Client {
	return &Client{server: s, host: h}
}



