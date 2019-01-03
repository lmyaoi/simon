package conf

import (
	"simon/jsonutil"
	"time"
)

var _default = Config{&jsonFormat{
	curVer,
	jsonutil.Duration{Duration: time.Second},
	9090,
	"localhost",
	8484,
	8484,
}}

type jsonFormat struct {
	Ver         string
	Interval    jsonutil.Duration
	VlcPort     int
	HostAddr    string
	HostPort    int
	HostingPort int
}

func New(interval time.Duration, hostUrl string, hostPort, hostingPort int) *Config {
	return &Config{
		&jsonFormat{
			curVer,
			jsonutil.Duration{Duration: interval},
			9090,
			hostUrl,
			hostPort,
			hostingPort,
		},
	}
}
