package local

import (
	"simon/conf"
	"simon/log"
	"simon/net"
	"simon/net/playback"
	"simon/ticker"
	"sync"
	"time"
)

type Client struct {
	server                       playback.Server
	host                         net.Host
	playbackControl, hostControl chan<- ticker.Signal
	playbackSignal, hostSignal   <-chan time.Time
	wg                           *sync.WaitGroup
}

func NewClient(s playback.Server, h net.Host, wg *sync.WaitGroup) *Client {
	c := &Client{server: s, host: h, wg: wg}
	c.playbackSignal, c.playbackControl = ticker.New(conf.Get().Interval.Duration)
	c.hostSignal, c.hostControl = ticker.New(conf.Get().Interval.Duration)
	go c.loop()
	return c
}

func (c *Client) loop() {
	for {
		select {
		case <-c.playbackSignal:
			if _, err := c.server.Status(); err != nil {
				log.Println(err)
			}
		case <-c.hostSignal:
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
	c.playbackControl <- ticker.On
	c.hostControl <- ticker.On
	c.server.On()
}

func (c *Client) Off() {
	c.playbackControl <- ticker.Off
	c.hostControl <- ticker.Off
	c.server.Off()
}

func (c *Client) Kill() {
	c.wg.Done()
	c.playbackControl <- ticker.Kill
	c.hostControl <- ticker.Kill
	c.server.Kill()
}
