package main

import (
	"flag"
	"time"
	"vsync/server"
)

func main() {
	interval := *flag.Duration("i", 1 * time.Second, "The interval with which to poll.")
	listen := *flag.Int("l", 4242, "The port on which to listen.")
	speak := *flag.Int("s", 8484, "The port on which to speak.")

	flag.Parse()
}
