package system

import (
	"context"
	"eff-gateway/gateway/proxy"
	"eff-gateway/setting"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type EffGateWay struct {
	server *http.Server

	locker  sync.RWMutex
	isClose bool

	closeChan chan byte
}

func Default() *EffGateWay {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", setting.Config.Server.Port),
		ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout),
		WriteTimeout: time.Duration(setting.Config.Server.WriteTimout),
		Handler:      initRouter(),
	}
	gateWay := &EffGateWay{
		server:    server,
		isClose:   false,
		closeChan: make(chan byte, 100),
	}
	return gateWay
}

func (g *EffGateWay) Run() {
	if g.server == nil {
		g.server = &http.Server{
			Addr:         fmt.Sprintf(":%d", setting.Config.Server.Port),
			ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout),
			WriteTimeout: time.Duration(setting.Config.Server.WriteTimout),
			Handler:      initRouter(),
		}
	}

	log.Println("Gateway service is running at port:", setting.Config.Server.Port)
	go g.osKill()
	g.runHTTPSv()
}

func (g *EffGateWay) runHTTPSv() {
	err := g.server.ListenAndServe()
	if err != nil {
		g.close()
		g.shutdown()
		log.Println("server is closing")
	}

}

func initRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.Forward(proxy.ProxyMap))
	return mux
}

func (g *EffGateWay) osKill() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
	g.close()
}

func (g *EffGateWay) close() {
	g.locker.Lock()
	defer g.locker.Unlock()

	if !g.isClose {
		g.isClose = true
		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(setting.Config.Server.ShutdownTimout)*time.Second)
		defer cancel()
		if err := g.server.Shutdown(ctx); err != nil {
			log.Println("shutting down: " + err.Error())
		}
		close(g.closeChan)
	}
}

func (g *EffGateWay) shutdown() {
	<-g.closeChan
	time.Sleep(time.Duration(setting.Config.Server.ShutdownTimout) * time.Second)
	log.Fatalln("Gateway service shutdown")
}
