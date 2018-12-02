package playback

import "time"

type Server interface {
	Connect() // connects to the playback server
	SetState(state State) // sets playback state
	Sync(Status) // syncs playback
	Status() Status // request current playback status
	Last() Status // request last requested playback status
}

type Status interface {
	State() State
	Pos() time.Time
	Created() time.Time
}

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
	return time.Duration(diff) >= 1 * time.Second
}
type State int

const (
	Stopped State = iota
	Playing
	Paused
)