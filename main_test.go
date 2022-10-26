package main

import (
	"context"
	"eff-gateway/discovery"
	"eff-gateway/gateway/proxy"
	"eff-gateway/gateway/proxy/types"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test1(t *testing.T) {

	//proxy.ProxyMap["/test"] = types.Proxy{
	//	IPAddr:        "http://127.0.0.1:9001",
	//	Prefix:        "/test",
	//	Upstream:      "/test",
	//	RewritePrefix: "",
	//}
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.Forward(proxy.ProxyMap))
	//http.HandleFunc("/", proxy.Forward(proxy.ProxyMap))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: mux,
	}
	log.Fatalln(server.ListenAndServe())
}

func TestEtcd(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.3.20:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer client.Close()
	fmt.Println("[INFO] connect to etcd success")
	//client.Put(context.Background(), "test", "data")
	response, err := client.Get(context.Background(), "test")
	fmt.Println(response.Kvs)
}

// 模拟服务
func TestServe(t *testing.T) {
	// 服务注册
	discovery.InitService()
	bytes, _ := json.Marshal(types.Proxy{
		IPAddr:        "http://127.0.0.1:9001",
		Prefix:        "/test",
		Upstream:      "/test",
		RewritePrefix: "",
	})
	discovery.PutSv("/test", string(bytes))
	fmt.Println("success")
	defer func() {
		discovery.DelSv("/test")
	}()

	type JsonResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	// handle all requests to your server using the proxy
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("收到")
		msg, _ := json.Marshal(JsonResult{Code: 200, Msg: "成功"})
		writer.Header().Set("content-type", "text/json")
		writer.Write(msg)
		return
	})

	log.Fatal(http.ListenAndServe(":9001", nil))
}

// 注册服务
func TestEtcdSend(t *testing.T) {
	discovery.InitEtcd()
	bytes, _ := json.Marshal(types.Proxy{
		IPAddr:        "http://127.0.0.1:9001",
		Prefix:        "/test",
		Upstream:      "/test",
		RewritePrefix: "",
	})
	discovery.EC.Put("/test", string(bytes))
	fmt.Println("success")

	time.Sleep(time.Second * 5)

	defer func() {
		discovery.EC.Del("/test")
	}()
}
