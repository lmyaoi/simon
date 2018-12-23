package net

import "vsync/net/playback"

type Host interface {
	Status() (playback.Status, error)
}
