package pref

import (
	"encoding/json"
	"errors"
	"os"
	"simon/jsonutil"
	"simon/net/playback/vlc/path"
	"strconv"
	"time"
)

var defaultPreference = Preference{&jsonFormat{
	"0",
	jsonutil.Duration{Duration: time.Second},
	path.Default,
	9090,
	"localhost",
	8484,
	8484,
}}

var p *Preference //singleton value

const name = ".vsync_prefs"

type Preference struct {
	data *jsonFormat
}

type jsonFormat struct {
	Ver         string
	Interval    jsonutil.Duration
	VlcPath     string
	VlcPort     int
	HostUrl     string
	HostPort    int
	HostingPort int
}

func (p *Preference) Interval() time.Duration {
	return p.data.Interval.Duration
}

func (p *Preference) VlcPath() string {
	return p.data.VlcPath
}

func (p *Preference) SetVlcPort(port int) {
	p.data.VlcPort = port
}

func (p *Preference) VlcPort() int {
	return p.data.VlcPort
}

func (p *Preference) VlcPortStr() string {
	return strconv.Itoa(p.data.VlcPort)
}

func (p *Preference) Host() (string, int) {
	return p.data.HostUrl, p.data.HostPort
}

func (p *Preference) HostUrl() string {
	return p.data.HostUrl
}

func (p *Preference) HostPort() int {
	return p.data.HostPort
}

func (p *Preference) HostStr() (string, string) {
	return p.data.HostUrl, strconv.Itoa(p.data.HostPort)
}

func (p *Preference) HostingPort() int {
	return p.data.HostingPort
}

func (p *Preference) HostingPortStr() string {
	return strconv.Itoa(p.data.HostingPort)
}

func New(interval time.Duration, vlcPath, hostUrl string, hostPort, hostingPort int) *Preference {
	return &Preference{
		&jsonFormat{
			"0",
			jsonutil.Duration{Duration: interval},
			vlcPath,
			9090,
			hostUrl,
			hostPort,
			hostingPort,
		},
	}
}

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
	err = json.NewEncoder(f).Encode(defaultPreference.data) // write default to new file
	if err != nil {
		return nil, err
	}
	_, _ = f.Seek(0, 0) // seek back to start of file
	return f, nil
}
