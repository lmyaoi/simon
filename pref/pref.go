package pref

import (
	"simon/jsonutil"
	"strconv"
	"time"
)

const curVer = "v1.0.0"

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

func New(interval time.Duration, vlcPath, hostUrl string, hostPort, hostingPort int) *Preference {
	return &Preference{
		&jsonFormat{
			curVer,
			jsonutil.Duration{Duration: interval},
			vlcPath,
			9090,
			hostUrl,
			hostPort,
			hostingPort,
		},
	}
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
