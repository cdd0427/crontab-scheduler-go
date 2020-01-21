package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort         int      `json:"api_port"`
	ApiReadTimeOut  int      `json:"api_read_time_out"`
	ApiWriteTimeOut int      `json:"api_write_time_out"`
	EtcdEndPoints   []string `json:"etcd_end_points"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
}

var (
	//singleton object
	G_config *Config
)

//load config from json file
func InitConig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//json deserialize
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf

	return
}
