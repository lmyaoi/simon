package conf

import (
	"simon/jsonutil"
	"time"
)

var _default = Config{
	curVer,
	jsonutil.Duration{Duration: time.Second},
	9090,
	"localhost",
	8484,
	8484,
}

type Config struct {
	Ver         string
	Interval    jsonutil.Duration
	VlcPort     int
	HostAddr    string
	HostPort    int
	HostingPort int
}

func New(interval time.Duration, hostUrl string, hostPort, hostingPort int) *Config {
	return &Config{
		curVer,
		jsonutil.Duration{Duration: interval},
		9090,
		hostUrl,
		hostPort,
		hostingPort,
	}
}
