package vlc

import (
	"vsync/log"
)

type Playlist struct {
	current int
	// todo: actual playlist
}

type Status struct {
	audioDelay int
	subtitleDelay int
	playlist Playlist
	paused bool
}

var Latest Status

func (*Status) String() string {
	return log.Sprintf("<INSERT DATA HERE>")
}
