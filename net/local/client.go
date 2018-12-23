package local

import (
	"sync"
	"time"
	"vsync/flags"
	"vsync/net"
	"vsync/log"
	"vsync/net/playback"
	"vsync/ticker"
)

type Client struct {
	server  playback.Server
	host    net.Host
	control chan<- ticker.Signal
	signal  <-chan time.Time
	wg      *sync.WaitGroup
}

func NewClient(s playback.Server, h net.Host, wg *sync.WaitGroup) *Client {
	c := &Client{server: s, host: h, wg: wg}
	c.signal, c.control = ticker.New(*flags.Interval)
	go c.loop()
	return c
}

func (c *Client) loop() {
	for {
		select {
		case <-c.signal:
			if _, err := c.server.Status(); err != nil {
				log.Println(err)
			}
			if stat, err := c.host.Status(); err != nil {
				log.Println(err)
			} else if err := c.server.Sync(stat); err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *Client) On() {
	c.wg.Add(1)
	c.control <- ticker.On
}

func (c *Client) Off() {
	c.control <- ticker.Off
}

func (c *Client) Kill() {
	c.wg.Done()
	c.control <- ticker.Kill
}
