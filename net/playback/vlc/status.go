package vlc

import (
	"encoding/json"
	"io"
	"simon/net/playback"
	"time"
)

type Status struct {
	*jsonStatus
}

type jsonStatus struct {
	State        playback.State
	Pos, Created time.Time
	Id           int
}

var DefaultStatus = Status{
	&jsonStatus{
		State:   playback.Stopped,
		Pos:     time.Unix(0, 0),
		Created: time.Now(),
		Id:      -1,
	},
}

func NewStatus(body io.Reader) *Status {
	s := parseJSON(body)
	return &Status{
		&jsonStatus{
			State:   parseState(s.State),
			Pos:     calcAccurateTime(s.Position, s.Length),
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

func (s *Status) SetState(state playback.State) {
	s.jsonStatus.State = state
}

func (s *Status) Pos() time.Time {
	return s.jsonStatus.Pos
}

func (s *Status) SetPos(p int64) {
	s.jsonStatus.Pos = time.Unix(p, 0)
}

func (s *Status) Created() time.Time {
	return s.jsonStatus.Created
}

func (s *Status) Id() int {
	return s.jsonStatus.Id
}

func (s *Status) SetId(id int) {
	s.jsonStatus.Id = id
}

func (s *Status) Copy() *Status {
	return &Status{
		&jsonStatus{
			State:   s.jsonStatus.State,
			Pos:     s.jsonStatus.Pos,
			Created: time.Now(),
			Id:      s.jsonStatus.Id,
		},
	}
}

func Verify(s playback.Status) *Status {
	t, ok := s.(*Status)
	if !ok {
		panic("Unexpected playback.Status implementation. Expected vlc.Status.")
	}
	return t
}

type vlcStatus struct {
	Length      int64
	Position    float64
	Time        int64
	State       string
	Currentplid int
}

func parseJSON(r io.Reader) *vlcStatus {
	s := &vlcStatus{}
	err := json.NewDecoder(r).Decode(s)
	if err != nil {
		panic(err)
	}
	return s
}

func parseState(str string) playback.State {
	switch str {
	case "stopped":
		return playback.Stopped
	case "playing":
		return playback.Playing
	case "paused":
		return playback.Paused
	default:
		panic("Impossible state string: " + str)
	}
}

func calcAccurateTime(p float64, length int64) time.Time {
	t := float64(length) * p * float64(time.Second)  // nanosecond time passed
	return time.Unix(0, int64(t)).Round(time.Second) // conversion to time rounded to second
}
