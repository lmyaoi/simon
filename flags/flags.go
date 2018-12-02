package flags

import (
	"flag"
	"time"
)

var (
	Interval = flag.Duration("i", 1 * time.Second, "The interval with which to poll.")
	HostUrl = flag.String("host-url", "localhost", "The url of the host.")
	HostPort = flag.Int("host-port", 8484, "The port on which the host listens.")
	VlcPort = flag.Int("vlc-port", 9090, "The port on which VLC listens.")
	Host = flag.Bool("host", false, "Whether this process is the host or the client.")
	VlcPwd = flag.String("pwd", "pwd", "The password used for the vlc http process.")
)
