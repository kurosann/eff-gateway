package setting

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	Config conf
)

func init() {
	Config = initConfig().checkConfig()
}

type conf struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
}

type ServerConfig struct {
	Port           int `yaml:"port"`
	ReadTimout     int `yaml:"readTimeout"`
	WriteTimout    int `yaml:"writeTimeout"`
	ShutdownTimout int `yaml:"shutdownTimeout"`
}

type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type Balance struct {
	Strategy      string `yaml:"strategy"`
	DefaultWeight string `yaml:"defaultWeight"`
}

func initConfig() conf {
	config, err := os.ReadFile("setting/setting.yaml")
	if err != nil {
		log.Println("读取配置文件失败:", err.Error())
		return conf{}
	}

	var c conf
	err = yaml.Unmarshal(config, &c)
	if err != nil {
		log.Println("配置文件序列化失败:", err.Error())
		return conf{}
	}
	return c
}

func (c conf) checkConfig() conf {
	if Config.Server.Port == 0 {
		Config.Server.Port = 8000
	}
	if Config.Server.ReadTimout == 0 {
		Config.Server.ReadTimout = 500
	}
	if Config.Server.WriteTimout == 0 {
		Config.Server.WriteTimout = 500
	}
	if Config.Server.ShutdownTimout == 0 {
		Config.Server.ShutdownTimout = 30
	}
	return c
}
