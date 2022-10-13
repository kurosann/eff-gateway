package types

type Server struct {
	ServerName string     `json:"server_name" yaml:"server_name"` // 服务名
	Location   []Location `json:"location" yaml:"location"`
}

type Location struct {
	LocationPath string `json:"location_path" yaml:"location_path"`
	ProxyPath    string `json:"proxy_path" yaml:"proxy_path"`
}
