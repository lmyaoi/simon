package conf

import (
	"simon/jsonutil"
	"time"
)

const curVer = "v1.0.0"

type Config struct {
	data *jsonFormat
}

type jsonFormat struct {
	Ver         string
	Interval    jsonutil.Duration
	VlcPath     string
	VlcPort     int
	HostAddr    string
	HostPort    int
	HostingPort int
}

func New(interval time.Duration, vlcPath, hostUrl string, hostPort, hostingPort int) *Config {
	return &Config{
		&jsonFormat{
			curVer,
			jsonutil.Duration{Duration: interval},
			vlcPath,
			9090,
			hostUrl,
			hostPort,
			hostingPort,
		},
	}
}
