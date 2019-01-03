// +build !darwin

package conf

import (
	"simon/jsonutil"
	"simon/net/playback/vlc/path"
	"time"
)

var _default = Config{&jsonFormat{
	curVer,
	jsonutil.Duration{Duration: time.Second},
	path.VlcDefault,
	9090,
	"localhost",
	8484,
	8484,
}}

type Config struct {
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
		curVer,
		jsonutil.Duration{Duration: interval},
		vlcPath,
		9090,
		hostUrl,
		hostPort,
		hostingPort,
	}
}
