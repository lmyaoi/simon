package preference

import (
	"encoding/json"
	"errors"
	"os"
	"time"
	"vsync/flags"
	"vsync/jsonutil"
)

var defaultPreference = Preference{&jsonFormat{"1", jsonutil.Duration{1 * time.Second}, flags.VlcDefault(), make([]string, 0)}}

const name = ".vsync_preferences"

type Preference struct {
	*jsonFormat
}

func (p *Preference) Interval() time.Duration {
	return p.jsonFormat.Interval.Duration
}

func (p *Preference) VlcPath() string {
	return p.jsonFormat.VlcPath
}

func (p *Preference) Hosts() []string {
	return p.jsonFormat.Hosts
}

type jsonFormat struct {
	Ver      string
	Interval jsonutil.Duration
	VlcPath  string
	Hosts    []string
}

func Get() (*Preference, error) {
	f, err := getFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseFile(f)
}

func parseFile(f *os.File) (*Preference, error) {
	format := &jsonFormat{}
	err := json.NewDecoder(f).Decode(format)
	if err != nil {
		return nil, err
	}
	if format.Ver != "1" {
		return nil, errors.New("invalid preference version")
	}
	return &Preference{format}, nil
}

func getFile() (*os.File, error) {
	path := os.Getenv("HOME") + "/" + name
	f, err := os.Open(path)
	if err == nil {
		return f, nil
	}
	if os.IsNotExist(err) {
		return createFile(path)
	}
	return nil, err // not a handleable error
}

func createFile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	err = json.NewEncoder(f).Encode(defaultPreference.jsonFormat) // write default to new file
	if err != nil {
		return nil, err
	}
	_, _ = f.Seek(0, 0) // seek back to start of file
	return f, nil
}
