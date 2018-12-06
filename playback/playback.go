package playback

import (
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

type Dummy int

const D Dummy = 0

func (Dummy) Connect() error             { return nil }      // connects to the playback server
func (Dummy) SetState(state State) error { return nil }      // sets playback state
func (Dummy) Sync(Status) error          { return nil }      // syncs playback
func (Dummy) Status() (Status, error)    { return nil, nil } // request current playback status
func (Dummy) Last() Status               { return nil }      // request last requested playback status

type Status interface {
	State() State
	Pos() time.Time
	Created() time.Time
	Marshal() []byte
}

type StatusUnmarshaler func(io.Reader) (Status, error)

func Now(s Status) time.Time {
	if s.State() == Paused {
		return s.Pos()
	} else if s.State() == Playing {
		return s.Pos().Add(time.Now().Sub(s.Created()))
	}
	return time.Time{}
}

func WorthSeeking(s1, s2 Status) bool {
	diff := Now(s1).UnixNano() - Now(s2).UnixNano()
	if diff < 0 {
		diff = -diff
	}
	return time.Duration(diff) >= 1*time.Second
}

type State int

const (
	Stopped State = iota
	Playing
	Paused
)
