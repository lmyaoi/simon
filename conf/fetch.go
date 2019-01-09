package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const name = ".simon_conf"

var c *Config //singleton value

func Initialize() {
	if c != nil {
		return
	}
	if err := Fetch(); err != nil {
		fmt.Println(err)
		c = &Config{}
		*c = _default
	}
}

func Fetch() error {
	f, err := getFile(os.O_RDONLY)
	if err != nil {
		return err
	}
	defer f.Close()
	c, err = parseFile(f)
	return err
}

func Get() *Config {
	return c
}

func Set(conf *Config) error {
	c = conf
	f, err := getFile(os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(conf)
}

func Save() error {
	return Set(c)
}

func parseFile(f *os.File) (*Config, error) {
	conf := &Config{}
	err := json.NewDecoder(f).Decode(conf)
	if err != nil {
		return nil, err
	}
	if conf.Ver != curVer {
		return nil, errors.New("invalid preference version")
	}
	return conf, nil
}

func getFile(flags int) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = createFile(); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, flags, 0666)
}

func createFile() error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(_default) // write default to new file
}
