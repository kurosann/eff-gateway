package main

import (
	"eff-gateway/gateway/juice"
)

func main() {
	gateWay := juice.Default()
	gateWay.Run()
}
