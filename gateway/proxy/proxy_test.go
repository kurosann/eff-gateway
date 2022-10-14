package proxy

import (
	"eff-gateway/balance"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
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

func TestProxyServer(t *testing.T) {

	type JsonResult struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	r := gin.Default()
	r.GET("/v2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":9001") // listen and serve on 0.0.0.0:8080

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
}
