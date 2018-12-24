package preference

import (
	"encoding/json"
	"errors"
	"os"
	"time"
	"vsync/flags"
	"vsync/jsonutil"
)

var defaultPreference = &Preference{&jsonFormat{
	"0",
	jsonutil.Duration{Duration: time.Second},
	flags.VlcDefault(),
	make([]string, 0),
}}

var p *Preference //singleton value

const name = ".vsync_prefs"

type Preference struct {
	*jsonFormat
}

type jsonFormat struct {
	Ver      string
	Interval jsonutil.Duration
	VlcPath  string
	Hosts    []string
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

func New(interval time.Duration, vlcPath string, hosts []string) *Preference {
	return &Preference{
		&jsonFormat{
			"0",
			jsonutil.Duration{Duration: interval},
			vlcPath,
			hosts,
		},
	}
}

func Get() (*Preference, error) {
	if p != nil {
		return p, nil
	}
	f, err := getFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p, err = parseFile(f)
	return p, err
}

func Set(preference *Preference) error {
	f, err := getFile()
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(preference)
	if err != nil {
		return err
	}
	p = preference
	return nil
}

func parseFile(f *os.File) (*Preference, error) {
	format := &jsonFormat{}
	err := json.NewDecoder(f).Decode(format)
	if err != nil {
		return nil, err
	}
	if format.Ver != "0" {
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
