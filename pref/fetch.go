package pref

import (
	"encoding/json"
	"errors"
	"os"
	"simon/jsonutil"
	"simon/net/playback/vlc/path"
	"time"
)

const name = ".simon_prefs"

var p *Preference //singleton value

var defaultPreference = Preference{&jsonFormat{
	"0",
	jsonutil.Duration{Duration: time.Second},
	path.Default,
	9090,
	"localhost",
	8484,
	8484,
}}

func Initialize() {
	if p != nil {
		return
	}
	if err := Fetch(); err != nil {
		println(err)
		p = &Preference{}
		*p = defaultPreference
	}
}

func Fetch() error {
	f, err := getFile()
	if err != nil {
		return err
	}
	defer f.Close()
	p, err = parseFile(f)
	return err
}

func Get() *Preference {
	return p
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
	if format.Ver != curVer {
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
	err = json.NewEncoder(f).Encode(defaultPreference.data) // write default to new file
	if err != nil {
		return nil, err
	}
	_, _ = f.Seek(0, 0) // seek back to start of file
	return f, nil
}
