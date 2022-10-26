package main

import (
	"eff-gateway/gateway/system"
)

func main() {
	gateWay := system.Default()
	gateWay.Run()
}
