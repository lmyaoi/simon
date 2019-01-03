package conf

import "encoding/json"

const curVer = "v1.0.0"

type Config struct {
	data *jsonFormat
}

func (conf *Config) String() string {
	b, _ := json.MarshalIndent(conf.data, "", "\t")
	return string(b)
}
