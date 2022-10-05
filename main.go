package main

import (
	"eff-gateway/gateway/system"
	"flag"
)

func main() {
	flag.Parse()
	gateWay := system.Default()
	gateWay.Run()
}
