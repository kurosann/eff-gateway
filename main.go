package main

import (
	"context"
	"flag"
	"fmt"
	"go-gateway/setting"
	"go-gateway/gateway/system"
	"log"
	"net/http"
	"time"
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", setting.Config.Server.Port),
		ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout),
		WriteTimeout: time.Duration(setting.Config.Server.WriteTimout),
		Handler:      system.InitRouter(),
	}
	go system.RunHTTPSv(server, cancel)
	go system.OsKill(cancel)
	<-ctx.Done()
	<-time.After(5 * time.Second)
	log.Fatalln("Gateway service stop")
}
