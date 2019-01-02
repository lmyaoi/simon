package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"simon/jsonutil"
	"simon/path"
	"time"
)

const name = ".simon_conf"

var c *Config //singleton value

var _default = Config{&jsonFormat{
	curVer,
	jsonutil.Duration{Duration: time.Second},
	path.VlcDefault,
	9090,
	"localhost",
	8484,
	8484,
}}

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
	return json.NewEncoder(f).Encode(conf.data)
}

func Save() error {
	return Set(c)
}

func parseFile(f *os.File) (*Config, error) {
	format := &jsonFormat{}
	err := json.NewDecoder(f).Decode(format)
	if err != nil {
		return nil, err
	}
	if format.Ver != curVer {
		return nil, errors.New("invalid preference version")
	}
	return &Config{format}, nil
}

func getFile(flags int) (*os.File, error) {
	filepath := fmt.Sprintf("%v%c%v", os.Getenv("HOME"), os.PathSeparator, name)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err = createFile(filepath); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filepath, flags, 0666)
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(_default.data) // write default to new file
}
