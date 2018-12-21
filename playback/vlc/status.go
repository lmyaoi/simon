package vlc

import (
	"encoding/json"
	"io"
	"time"
	"vsync/playback"
)

type Status struct {
	*jsonStatus
}

type jsonStatus struct {
	State        playback.State
	Pos, Created time.Time
	Id           int
}

func NewStatus(body io.Reader) *Status {
	s := json2struct(body)
	return &Status{
		&jsonStatus{
			State:   string2state(s.State),
			Pos:     time.Unix(calcAccurateTime(s.Position, s.Length)),
			Created: time.Now(),
			Id:      s.Currentplid,
		},
	}
}

func Unmarshal(body io.Reader) (playback.Status, error) {
	js := &jsonStatus{}
	if err := json.NewDecoder(body).Decode(js); err != nil {
		return nil, err
	}
	return &Status{js}, nil
}

func (s *Status) Marshal() []byte {
	data, err := json.Marshal(s.jsonStatus)
	if err != nil {
		panic(err)
	}
	return data
}

func (s *Status) State() playback.State {
	return s.jsonStatus.State
}

func (s *Status) Pos() time.Time {
	return s.jsonStatus.Pos
}

func (s *Status) Created() time.Time {
	return s.jsonStatus.Created
}

func (s *Status) Id() int {
	return s.jsonStatus.Id
}

func verify(s playback.Status) *Status {
	t, ok := s.(*Status)
	if !ok {
		panic("Unexpected playback.Status implementation. Expected vlc.Status.")
	}
	return t
}

func json2struct(r io.Reader) *struct {
	Length      int64
	Position    float64
	State       string
	Currentplid int
} {
	s := &struct {
		Length      int64
		Position    float64
		State       string
		Currentplid int
	}{}
	err := json.NewDecoder(r).Decode(s)
	if err != nil {
		panic(err)
	}
	return s
}

func string2state(str string) playback.State {
	switch str {
	case playback.Stopped.String():
		return playback.Stopped
	case playback.Playing.String():
		return playback.Playing
	case playback.Paused.String():
		return playback.Paused
	default:
		panic("Impossible state string")
	}
}

func calcAccurateTime(p float64, length int64) (int64, int64) {
	t := float64(length) * p
	s := int64(t)
	ns := int64((t - float64(s)) * float64(time.Second))
	return s, ns
}
