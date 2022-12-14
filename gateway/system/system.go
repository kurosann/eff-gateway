package system

import (
	"context"
	"eff-gateway/discovery"
	"eff-gateway/gateway/proxy"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/setting"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
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
		ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout) * time.Microsecond,
		WriteTimeout: time.Duration(setting.Config.Server.WriteTimout) * time.Microsecond,
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
			ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout) * time.Microsecond,
			WriteTimeout: time.Duration(setting.Config.Server.WriteTimout) * time.Microsecond,
			Handler:      initRouter(),
		}
	}
	g.init()
	log.Println("Gateway service is running at port:", setting.Config.Server.Port)

	go g.osKill()
	g.runHTTPSv()
}

func (g *EffGateWay) init() {
	discovery.InitService()
	// etcd 服务发现
	go discovery.KeepAlive("/test", func(ev *clientv3.Event) bool {
		fmt.Printf("Servise update type:%s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		if ev.Type == mvccpb.PUT {
			v := &types.Proxy{}
			err := json.Unmarshal(ev.Kv.Value, v)
			if err != nil {
				return true
			}
			proxy.ProxyMap[string(ev.Kv.Key)] = *v
			return true
		}
		if ev.Type == mvccpb.DELETE {
			delete(proxy.ProxyMap, string(ev.Kv.Key))
		}
		return true
	})
	proxy.InitProxy()
}

func (g *EffGateWay) runHTTPSv() {
	err := g.server.ListenAndServe()
	if err != nil {
		g.close()
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
		g.shutdown()
	}
}

func (g *EffGateWay) shutdown() {
	time.Sleep(time.Duration(setting.Config.Server.ShutdownTimout) * time.Second)
	discovery.EC.Close()
	log.Fatalln("Gateway service shutdown")
}
