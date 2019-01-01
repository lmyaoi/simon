package jsonutil

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration{dur}
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
