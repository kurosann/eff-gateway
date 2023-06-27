package juice

import (
	"context"
	"eff-gateway/discovery"
	"eff-gateway/gateway/proxy"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/glog"
	"eff-gateway/setting"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type GateWay struct {
	server *http.Server

	locker  sync.RWMutex
	isClose bool

	closeChan chan byte
}

func Default() *GateWay {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", setting.Config.Server.Port),
		ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout) * time.Microsecond,
		WriteTimeout: time.Duration(setting.Config.Server.WriteTimout) * time.Microsecond,
		Handler:      initRouter(),
	}
	gateWay := &GateWay{
		server:    server,
		isClose:   false,
		closeChan: make(chan byte, 100),
	}
	return gateWay
}

func (g *GateWay) Run() {
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

	go g.runHTTPSv()
	g.osKill()

}

func (g *GateWay) init() {
	err := discovery.InitService()
	err = discovery.PutSv("/juice/GateWay", setting.Config.Server.Ip)
	if err != nil {
		glog.ErrorLog.Println(err.Error())
		glog.ErrorLog.Fatalln("etcd service is run fail!")
	}
	glog.InfoLog.Println("[INFO] connect to etcd success")
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

func (g *GateWay) runHTTPSv() {
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

func (g *GateWay) osKill() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
	g.close()
}

func (g *GateWay) close() {
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

func (g *GateWay) shutdown() {
	time.Sleep(time.Duration(setting.Config.Server.ShutdownTimout) * time.Second)
	discovery.EC.Close()
	log.Fatalln("Gateway service shutdown")
}
