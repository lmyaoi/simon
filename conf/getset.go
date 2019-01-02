package conf

import (
	"time"
)

func (p *Config) Interval() time.Duration {
	return p.data.Interval.Duration
}
func (p *Config) SetInterval(interval time.Duration) {
	p.data.Interval.Duration = interval
}

func (p *Config) VlcPath() string {
	return p.data.VlcPath
}

func (p *Config) SetVlcPath(path string) {
	p.data.VlcPath = path
}

func (p *Config) SetVlcPort(port int) {
	p.data.VlcPort = port
}

func (p *Config) VlcPort() int {
	return p.data.VlcPort
}

func (p *Config) HostAddr() string {
	return p.data.HostAddr
}

func (p *Config) SetHostAddr(addr string) {
	p.data.HostAddr = addr
}

func (p *Config) HostPort() int {
	return p.data.HostPort
}

func (p *Config) SetHostPort(port int) {
	p.data.HostPort = port
}

func (p *Config) HostingPort() int {
	return p.data.HostingPort
}
