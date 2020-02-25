package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndPoints         []string `json:"etcd_end_points"`
	EtcdDialTimeout       int      `json:"etcd_dial_timeout"`
	MongodbUri            string   `json:"mongodb_uri"`
	MongodbConnectTimeout int      `json:"mongodb_connect_timeout"`
	JobLogBatchSize       int      `json:"job_log_batch_size"`
	JobLogCommitTimeout   int      `json:"job_log_commit_timeout"`
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
