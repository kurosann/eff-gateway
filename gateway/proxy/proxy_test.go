package proxy

import (
	"eff-gateway/balance"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	bl := balance.GlobalStrategy
	bl.AddStrategy("localhost1")
	bl.GetServer("localhost1").Impl.Add("http://127.0.0.1:9001", 1)

	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("http://127.0.0.1:9001", "localhost1")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func TestProxyServer(t *testing.T) {

	type JsonResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("收到")
		msg, _ := json.Marshal(JsonResult{Code: 400, Msg: "验证失败"})
		writer.Header().Set("content-type", "text/json")
		writer.Write(msg)
		return
	})

	log.Fatal(http.ListenAndServe(":9001", nil))
}
