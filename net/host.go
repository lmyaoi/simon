package net

import "simon/net/playback"

type Host interface {
	Status() (playback.Status, error)
}
