package vlc

import (
	"encoding/json"
	"io"
	"time"
	"vsync/playback"
)

type Status struct {
	state playback.State
	pos, created time.Time
	id int
}

func json2struct (r io.Reader) *struct {Length int64; Position float64; State string; Currentplid int} {
	s := &struct {Length int64; Position float64; State string; Currentplid int}{}
	err := json.NewDecoder(r).Decode(s)
	if err != nil { panic(err) }
	return s
}

func string2state (str string) playback.State {
	switch str {
	case "stopped": return playback.Stopped
	case "playing": return playback.Playing
	case "paused": return playback.Paused
	default: panic("Impossible state string")
	}
}

func NewStatus(body io.Reader) *Status {
	obj := json2struct(body)
	return &Status{
		state: string2state(obj.State),
		pos: time.Unix(toAccurateTime(obj.Position, obj.Length)),
		created: time.Now(),
		id: obj.Currentplid,
	}
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