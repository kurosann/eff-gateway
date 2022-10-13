package proxy

import (
	"eff-gateway/balance"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	LoadProxyUrlMap = make(map[string]LoadProxy)
)

type LoadProxy struct {
	ServerName string
	ProxyPath  string
}

// HostReverseProxyV1
// 构建基本的代理方法
// w 响应写入
// r http请求
func HostReverseProxyV1(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		io.WriteString(w, "Request path Error")
		return
	}
	var loadProxy = matchUrl(r.RemoteAddr, r.URL.RequestURI())

	bl := balance.GlobalStrategy
	proxyHost := bl.GetServer(loadProxy.ServerName).Impl.GetNode(loadProxy.ServerName)
	proxyHost += loadProxy.ProxyPath
	// parse the url
	reverseProxy, err := NewProxy(proxyHost)

	if err != nil {
		panic(err)
	}
	reverseProxy.ServeHTTP(w, r)
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	// 解析目标地址
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return errors.New("response body is invalid")
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

// matchUrl
// 匹配路径，返回代理的服务信息
func matchUrl(hostUrl, allPath string) LoadProxy {
	lp := LoadProxy{}
	if v, ok := LoadProxyUrlMap[allPath]; ok {
		return v
	}
	// 匹配地址如果含
	// 如：
	// k = http://127.0.0.1:8001
	// httpUrl = http://127.0.0.1:8001/admin-app/v1/ap/
	cache := LocalProxyCache[hostUrl]
	// 查找服务里的定位
	for _, server := range cache.ServerMap {
		for _, location := range server.Location {
			// location = admin-app
			// httpUrl = http://127.0.0.1:8001/admin-app/v1/ap/
			if strings.Contains(allPath, location.LocationPath) {
				lp.ServerName = server.ServerName
				lp.ProxyPath = location.ProxyPath
				if location.ProxyPath == "" || location.ProxyPath == "/" {
					lp.ProxyPath = strings.Split(allPath, hostUrl)[1]
				}
				LoadProxyUrlMap[allPath] = lp
				return lp
			}
		}
	}
	return LoadProxy{}
}
