package types

type Response struct {
	Success bool     `json:"success" yaml:"success"`
	Status  string   `json:"status" yaml:"status"`
	Server  []Server `json:"server" yaml:"server"`
	Data    []Proxy  `json:"data" yaml:"data"`
}
