package types

type GlobalHttp struct {
	HttpList []HttpConfig `json:"http_list" yaml:"http_list"`
}

type HttpConfig struct {
	Scheme           string   `json:"scheme" yaml:"scheme"`
	IpAddr           string   `json:"ip_addr" yaml:"ip_Addr"`
	Domain           string   `json:"domain" yaml:"domain"`
	Listen           string   `json:"listen" yaml:"listen"` // 监听端口
	DefaultType      string   `json:"default_type" yaml:"default_type"`
	SendFile         string   `json:"sendfile" yaml:"sendfile"`
	KeepaliveTimeout string   `json:"keepalive_timeout" yaml:"keepalive_timeout"`
	Server           []Server `json:"server" yaml:"server"`
}
