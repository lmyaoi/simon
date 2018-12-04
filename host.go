package main

import "vsync/playback"

type Host interface {
	Status() playback.Status
}