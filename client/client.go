package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
	"vsync/flags"
	"vsync/log"
	"vsync/ticker"
	"vsync/vlc"
)

type Client struct {
	ticker chan<- ticker.Signal
	c *http.Client
	url string
	latest vlc.Status
	wg *sync.WaitGroup
}

func New(wg *sync.WaitGroup) *Client {
	url := fmt.Sprintf("http://%v:%v", *flags.Url, *flags.Port)
	client := &http.Client{}
	ticks, control := ticker.New(*flags.Interval)
	c := &Client{control, client, url, vlc.Status{}, wg}
	go c.loop(ticks)
	return c
}

func (c *Client) Connect(retries int) {
	for retries >= 0 {
		r, err := c.c.Do(c.makeRequest("connect"))
		if err != nil || r.StatusCode != http.StatusOK {
			fmt.Printf("%v Failed to connect. Retrying...\n", log.Now())
			retries--
			time.Sleep(1 * time.Second)
		} else {
			fmt.Printf("%v Successfuly connected to host: %v\n", log.Now(), c.url)
			c.ticker <- ticker.On
			c.wg.Add(1)
			return
		}
	}
	c.wg.Done()
	panic(fmt.Sprintf("%v Failed to connect to host: %v\n", time.Now(), c.url))
}

func (c *Client) update() {
	r, err := c.c.Do(c.makeRequest("update"))
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r.Body)
}

func (c *Client) loop(ticks <-chan time.Time) {
	for {
		<-ticks
		c.update()
	}
}


func (c *Client) makeRequest(s string) *http.Request {
		url := fmt.Sprintf("%v/%v", c.url, s)
		return must(http.NewRequest("GET", url, nil))
}

func must(r *http.Request, err error) *http.Request {
	if err != nil {
		panic(err)
	}
	return r
}
