package flags

import (
	"flag"
	"time"
)

var (
	Interval = flag.Duration("i", 1 * time.Second, "The interval with which to poll.")
	Port = flag.Int("p", 8484, "The port on which to listen/talk.")
	Url = flag.String("u", "", "The url of the host (only necessary if client)")
	Host = flag.Bool("host", false, "Whether this process is the host or the client.")
)
