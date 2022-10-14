package main

import (
	"eff-gateway/gateway/proxy"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/gateway/system"
)

func main() {
	proxy.ProxyMap["/api/v1/add/"] = types.Proxy{IPAddr: "127.0.0.1:9001",
		Prefix:        "",
		Upstream:      "/test",
		RewritePrefix: ""}
	gateWay := system.Default()
	gateWay.Run()
}
