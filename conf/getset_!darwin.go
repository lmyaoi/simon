// +build !darwin

package conf

func (p *Config) VlcPath() string {
	return p.data.VlcPath
}

func (p *Config) SetVlcPath(path string) {
	p.data.VlcPath = path
}
