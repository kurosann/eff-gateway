package setting

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Config conf
)

func init() {
	getConfig(&Config)
	checkConfig()
}

type conf struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
}

type ServerConfig struct {
	Port        int `json:"port"`
	ReadTimout  int `json:"readTimout"`
	WriteTimout int `json:"writeTimout"`
}

type LogConfig struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

func getConfig(conf *conf) {
	config, err := ioutil.ReadFile("../setting/setting.yml")
	if err != nil {
		log.Println("读取配置文件失败")
	}
	err = yaml.Unmarshal(config, &conf)
	if err != nil {
		log.Println("配置文件序列化失败")
	}
}

func checkConfig() {
	if Config.Server.Port == 0 {
		Config.Server.Port = 8000
	}
	if Config.Server.ReadTimout == 0 {
		Config.Server.ReadTimout = 500
	}
	if Config.Server.WriteTimout == 0 {
		Config.Server.WriteTimout = 500
	}
}
