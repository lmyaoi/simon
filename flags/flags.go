package flags

import (
	"flag"
	"runtime"
	"time"
)

const (
	VLC_DEFAULT_WIN = `C:\Program Files\VideoLAN\VLC\vlc.exe`
	VLC_DEFAULT_OSX = `/Applications/VLC.app`
)

var (
	Interval = flag.Duration("i", 1*time.Second, "The interval at which to poll.")
	HostUrl  = flag.String("host-url", "localhost", "The url/ip of the host.")
	HostPort = flag.Int("host-port", 8484, "The port on which the host listens.")
	VlcPort  = flag.Int("vlc-port", 9090, "The port on which VLC listens.")
	Host     = flag.Bool("host", false, "Set when host.")
	vlcPath  = flag.String("vlc", "", "The path to the vlc executable.")
)

func vlcpath() string {
	switch runtime.GOOS {
	case `windows`:
		return VLC_DEFAULT_WIN
	case `darwin`:
		return VLC_DEFAULT_OSX
	default:
		return ""
	}
}

func Vlc() string {
	if runtime.GOOS == `darwin` {
		return *vlcPath + `/Contents/MacOS/VLC`
	}
	return *vlcPath
}
