package vlc

import (
	"encoding/json"
	"io"
	"time"
	"vsync/playback"
)

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

func Marshal(status Status) []byte {
	data, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	return data
}

func Unmarshal(r io.Reader) *Status {
	s := json2struct(r)
	return &Status{
		state:   string2state(s.State),
		pos:     time.Unix(toAccurateTime(s.Position, s.Length)),
		created: time.Now(),
		id:      s.Currentplid,
	}
}
