package flags

import (
	"flag"
	"runtime"
	"time"
)

const (
	VlcPathWin = `C:\Program Files\VideoLAN\VLC\vlc.exe`
	VlcPathOSX = `/Applications/VLC.app`
)

var (
	Interval = flag.Duration("i", time.Second, "The interval at which to poll.")
	HostUrl  = flag.String("host-url", "localhost", "The url/ip of the host.")
	HostPort = flag.Int("host-port", 8484, "The port on which the host listens.")
	VlcPort  = flag.Int("vlc-port", 9090, "The port on which VLC listens.")
	Host     = flag.Bool("host", false, "Set when host.")
	VlcPath      = flag.String("vlc-path", _vlcPath(), "The path to the vlc executable.")
)

func _vlcPath() string {
	switch runtime.GOOS {
	case `windows`:
		return VlcPathWin
	case `darwin`:
		return VlcPathOSX
	default:
		return ""
	}
}

func Vlc() string {
	if runtime.GOOS == `darwin` {
		return *VlcPath + `/Contents/MacOS/VLC`
	}
	return *VlcPath
}
