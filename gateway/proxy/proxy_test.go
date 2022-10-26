package proxy

import (
	"eff-gateway/balance"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/setting"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go func() {
		bl := balance.GlobalStrategy
		bl.AddStrategy("localhost1")
		bl.GetServer("localhost1").Impl.Add("http://127.0.0.1:9001", 1)
	}()
	// initialize a reverse proxy and pass the actual backend server url here
	//proxy, err := NewProxy("http://127.0.0.1:9001", "localhost1")
	//if err != nil {
	//	panic(err)
	//}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", HostReverseProxyV1)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//func TestProxyServer(t *testing.T) {
//
//	type JsonResult struct {
//		Code int    `json:"code"`
//		Msg  string `json:"msg"`
//	}
//	r := gin.Default()
//	r.GET("/v2", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "pong",
//		})
//	})
//
//	r.Run(":9001") // listen and serve on 0.0.0.0:8080

// handle all requests to your server using the proxy
//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println("收到")
//	msg, _ := json.Marshal(JsonResult{Code: 400, Msg: "验证失败"})
//	writer.Header().Set("content-type", "text/json")
//	writer.Write(msg)
//	return
//})
// handle all requests to your server using the proxy
//http.HandleFunc("/localhost1/v1", func(writer http.ResponseWriter, request *http.Request) {
//	fmt.Println("收到v1")
//	msg, _ := json.Marshal(JsonResult{Code: 400, Msg: "验证失败"})
//	writer.Header().Set("content-type", "text/json")
//	writer.Write(msg)
//	return
//})
//
//log.Fatal(http.ListenAndServe(":9001", nil))
//}

func TestServe(t *testing.T) {

	type JsonResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("收到")
		msg, _ := json.Marshal(JsonResult{Code: 400, Msg: "验证失败"})
		writer.Header().Set("content-type", "text/json")
		writer.Write(msg)
		return
	})

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func TestProxyReq(t *testing.T) {
	ProxyMap["/test"] = types.Proxy{
		IPAddr:        "http://127.0.0.1:9001",
		Prefix:        "/test",
		Upstream:      "/test",
		RewritePrefix: "",
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", Forward(ProxyMap))
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8081),
		ReadTimeout:  time.Duration(setting.Config.Server.ReadTimout),
		WriteTimeout: time.Duration(setting.Config.Server.WriteTimout),
		Handler:      mux,
	}
	log.Fatalln(server.ListenAndServe())
}
