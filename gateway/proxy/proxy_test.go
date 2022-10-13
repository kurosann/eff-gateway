package proxy

import (
	"log"
	"net/http"
	"testing"
)

func TestProxyReqs(t *testing.T) {
	// initialize a reverse proxy and pass the actual backend server url here
	proxy, err := NewProxy("http://my-api-server.com")
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
