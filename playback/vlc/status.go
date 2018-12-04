package vlc

import (
	"io"
	"time"
	"vsync/playback"
)

type Status struct {
	state playback.State
	pos, created time.Time
	id int
}

func NewStatus(body io.Reader) *Status {
	return Unmarshal(body)
}

func toAccurateTime(p float64, length int64) (int64, int64) {
	t := float64(length) * p
	s := int64(t)
	ns := int64((t - float64(s)) * float64(time.Second))
	return s, ns
}

func (s *Status) State() playback.State {
	return s.state
}

func (s *Status) Pos() time.Time {
	return s.pos
}

func (s *Status) Created() time.Time {
	return s.created
}

func (s *Status) Id() int {
	return s.id
}

func verify(s playback.Status) *Status {
	t, _ := s.(*Status)
	return t
}