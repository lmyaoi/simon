package conf

import "encoding/json"

const curVer = "v1.0.0"


func (conf *Config) String() string {
	b, _ := json.MarshalIndent(conf, "", "\t")
	return string(b)
}
