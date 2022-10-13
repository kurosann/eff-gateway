package types

type Proxy struct {
	ServiceName   string `json:"serviceName" yaml:"serviceName"`      // 服务名
	Domain        string `json:"domain" yaml:"domain"`                // 域名
	IPAddr        string `json:"ip_addr" yaml:"ip_addr"`              // IP地址
	Remark        string `json:"remark" yaml:"remark"`                // 备注
	Prefix        string `json:"prefix" yaml:"prefix"`                // 前缀
	Upstream      string `json:"upstream" yaml:"upstream"`            // 上游部分
	RewritePrefix string `json:"rewrite_prefix" yaml:"rewritePrefix"` // 重写前缀
}
