package main

import (
	"eff-gateway/gateway/proxy"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/gateway/system"
)

func main() {
	proxy.ProxyMap["/test"] = types.Proxy{
		IPAddr:        "http://127.0.0.1:9001",
		Prefix:        "/test",
		Upstream:      "/test",
		RewritePrefix: "",
	}
	gateWay := system.Default()
	gateWay.Run()
}
