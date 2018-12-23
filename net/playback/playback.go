package playback

import (
	"encoding/json"
	"io"
	"time"
)

type Server interface {
	Connect() error             // connects to the playback server
	SetState(state State) error // sets playback state
	Sync(Status) error          // syncs playback
	Status() (Status, error)    // request current playback status
	Last() Status               // request last requested playback status
}

type dummy int

const Dummy dummy = 0

func (dummy) Connect() error             { return nil }      // connects to the playback server
func (dummy) SetState(state State) error { return nil }      // sets playback state
func (dummy) Sync(Status) error          { return nil }      // syncs playback
func (dummy) Status() (Status, error)    { return nil, nil } // request current playback status
func (dummy) Last() Status               { return nil }      // request last requested playback status

type Status interface {
	State() State
	Pos() time.Time
	Created() time.Time
	Marshal() []byte
}

type StatusUnmarshaler func(io.Reader) (Status, error)

func Now(s Status) time.Time {
	switch s.State() {
	case Paused:
		return s.Pos()
	case Playing:
		return s.Pos().Add(time.Since(s.Created()))
	default:
		return time.Unix(0, 0)
	}
}

func WorthSeeking(s1, s2 Status) bool {
	diff := Now(s1).Sub(Now(s2)).Round(time.Second)
	if diff < 0 {
		diff = -diff
	}
	return diff > time.Second
}

//go:generate stringer -type=State
type State int

const (
	Stopped State = iota
	Playing
	Paused
)

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(s))
}

func (s *State) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*s = State(i)
	return nil
}
