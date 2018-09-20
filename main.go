package main

import (
	"flag"
	"sync"
	"vsync/client"
	"vsync/flags"
	"vsync/host"
)

func main() {

	flag.Parse()
	wg := &sync.WaitGroup{}
	if *flags.Host {
		h := host.New()
		h.Listen()
	} else {
		c := client.New(wg)
		c.Connect(10)
	}
	wg.Wait()
}
