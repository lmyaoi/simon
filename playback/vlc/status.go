package vlc

import (
	"encoding/json"
	"io"
	"time"
	"vsync/playback"
)

type Status struct {
	state        playback.State
	pos, created time.Time
	id           int
}

func NewStatus(body io.Reader) *Status {
	s := json2struct(body)
	return &Status{
		state:   string2state(s.State),
		pos:     time.Unix(calcAccurateTime(s.Position, s.Length)),
		created: time.Now(),
		id:      s.Currentplid,
	}
}

type marshalable struct {
	State        int
	Pos, Created int64
	Id           int
}

func (s *Status) toMS() *marshalable {
	return &marshalable{
		State:   int(s.state),
		Pos:     s.pos.UnixNano(),
		Created: s.created.UnixNano(),
		Id:      s.id,
	}
}

func (ms *marshalable) toS() *Status {
	return &Status{
		state:   playback.State(ms.State),
		pos:     time.Unix(0, ms.Pos),
		created: time.Unix(0, ms.Created),
		id:      ms.Id,
	}
}

func Unmarshal(body io.Reader) (playback.Status, error) {
	ms := &marshalable{}
	if err := json.NewDecoder(body).Decode(ms); err != nil {
		return nil, err
	}
	return ms.toS(), nil
}

func (s *Status) Marshal() []byte {
	data, err := json.Marshal(s.toMS())
	if err != nil {
		panic(err)
	}
	return data
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
	case "stopped":
		return playback.Stopped
	case "playing":
		return playback.Playing
	case "paused":
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
