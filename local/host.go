package local

import "vsync/playback"

type Host struct {
	server playback.Server
}

func (h *Host) Status() playback.Status {
	return h.server.Status()
}

