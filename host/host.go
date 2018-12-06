package host

import "vsync/playback"

type Host interface {
	Status() (playback.Status, error)
}
