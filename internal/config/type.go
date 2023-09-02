package config

import "github.com/gin-gonic/gin"

type Config struct {
	Debug      *bool `comment:"simple ip-kvm will look for a configuration file in the following order:\n - ./.simple-ip-kvm.yaml\n - ~/.simple-ip-kvm.yaml\nor you Or you can specify it through the parameter --config "`
	SerialPort string
	Web        struct {
		Addr     string
		BasePath string
		Title    string
		Auth     gin.Accounts `comment:"key/value for user/pass list for http basic authorize, empty to disable"`
	}
	Video struct {
		Type         VideoStreamerType `comment:"video streamer type, [external]"`
		StreamUrl    string
		StreamPath   string
		SnapshotPath string

		Streamer struct {
			Port int
		}
	}
	Proxy struct {
		Enable  bool
		Path    string
		Rewrite bool `comment:"rewrite path to /"`
	} `comment:"reverse proxy to video streamer"`
}

var cfg Config

func GetConfig() *Config {
	return &cfg
}
