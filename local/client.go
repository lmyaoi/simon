package local

import (
	"sync"
	"time"
	"vsync/flags"
	"vsync/host"
	"vsync/playback"
	"vsync/ticker"
)

type Client struct {
	server playback.Server
	host   host.Host
	control chan<- ticker.Signal
	signal <-chan time.Time
	wg *sync.WaitGroup
}

func NewClient(s playback.Server, h host.Host, wg *sync.WaitGroup) *Client {
	c := &Client{server: s, host: h, wg: wg}
	c.signal, c.control = ticker.New(*flags.Interval)
	go c.loop()
	return c
}

func (c *Client) loop() {
	for {
		select {
		case <-c.signal:
			c.server.Status()
			c.server.Sync(c.host.Status())
		}
	}
}

func (c *Client) On() {
	c.wg.Add(1)
	c.control<-ticker.On
}

func (c *Client) Off() {
	c.control<-ticker.Off
}

func (c *Client) Kill() {
	c.wg.Done()
	c.control<-ticker.Kill
}