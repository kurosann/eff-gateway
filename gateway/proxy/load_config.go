package proxy

import (
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/glog"
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// LocalProxyCache
// 本地代理缓存
var (
	LocalProxyCache   = make(map[string]LocalCache)
	LocalJsonFilePath = "config_local.json"
	LocalYamlFilePath = "config_local.yml"
)

type LocalCache struct {
	GlobalHttp types.HttpConfig
	ServerMap  map[string]types.Server
}

func init() {
	InitConfig()
}

func GetLocalProxyCache(name string) LocalCache {
	if v, k := LocalProxyCache[name]; k {
		return v
	} else {
		InitConfig()
		return LocalProxyCache[name]
	}
}

func InitConfig() {
	var respond types.GlobalHttp
	err := loadYaml(&respond, "")
	if err != nil {
		panic(err)
	}
	if len(respond.HttpList) == 0 {
		loadJson(&respond, "")
	}
	storeLocalConfig(respond)
	checkConfig()
}

// LoadJson
// 解析本地配置
func loadJson(respond interface{}, path string) {
	if path == "" {
		// 启用默认配置
		path = LocalJsonFilePath
	}
	// 解析json文件
	file, err := os.Open(path)
	if err != nil {
		glog.ErrorLog.Println("err:", err)
	}
	// 创建json解码器
	decoder := json.NewDecoder(file)
	err = decoder.Decode(respond)
	if err != nil {
		fmt.Println("LoadProxyListFromFile failed", err.Error())
	}
}

// LoadYaml
// 解析本地配置
func loadYaml(respond interface{}, path string) error {
	if path == "" {
		// 启用默认配置
		path = LocalYamlFilePath
	}
	// 创建json解码器
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, respond)
	if err != nil {
		return fmt.Errorf("in file %q: %v", path, err)
	}
	return nil

}

// 检测文件配置
func checkConfig() {
	fmt.Println("LocalProxyCache 检查")
	for key, value := range LocalProxyCache {
		fmt.Printf("LocalProxyCache key %s \n", key)
		for k, v := range value.ServerMap {
			fmt.Printf("[key:%s,v:%s ]\n", k, v)
		}

	}
}

// StoreLocalJson
// 遍历数据结构，缓存数据
func storeLocalConfig(response types.GlobalHttp) {
	// 遍历全部的配置文件的http配置
	for _, config := range response.HttpList {
		// 遍历http配置里的server
		var cacheMap = LocalCache{}
		cacheMap.ServerMap = map[string]types.Server{}
		for i, v := range config.Server {
			cacheMap.ServerMap[v.ServerName] = config.Server[i]
		}
		cacheMap.GlobalHttp = config
		if config.IpAddr != "" {
			LocalProxyCache[config.IpAddr+":"+config.Listen] = cacheMap
		}

		if config.Domain != "" {
			LocalProxyCache[config.Domain+":"+config.Listen] = cacheMap
		}

	}
}
