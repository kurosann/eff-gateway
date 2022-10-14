package proxy

import (
	"eff-gateway/balance"
	"eff-gateway/glog"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
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
	var loadProxy = matchUrl(r.Host, r.URL.RequestURI())
	bl := balance.GlobalStrategy
	proxyHost := bl.GetServer(loadProxy.ServerName).Impl.GetNode(loadProxy.ServerName)
	proxyHost += loadProxy.ProxyPath
	// parse the url
	reverseProxy, err := NewProxy(proxyHost, loadProxy.ServerName)

	if err != nil {
		panic(err)
	}
	reverseProxy.ServeHTTP(w, r)

}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost, serverName string) (*httputil.ReverseProxy, error) {
	// 解析目标地址
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	glog.InfoLog.Printf("RequestURI %s Path %s ", targetUrl.RequestURI(), targetUrl.Path)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req, targetUrl)
	}

	proxy.ModifyResponse = modifyResponse(targetUrl.Scheme+"://"+targetUrl.Host, serverName)
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request, u *url.URL) {
	// 自定义请求前的操作
	// 如下
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.URL = u
	req.RequestURI = u.RequestURI()

}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

func modifyResponse(targetUrl, serverName string) func(*http.Response) error {
	startTime := time.Now().UnixMilli()
	return func(resp *http.Response) error {
		bl := balance.GlobalStrategy
		bl.GetServer(serverName).Impl.AddReqs(targetUrl, int(time.Now().UnixMilli()-startTime))
		//return errors.New("response body is invalid")
		return nil
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
	cache := GetLocalProxyCache(hostUrl)
	// 查找服务里的定位
	for _, server := range cache.ServerMap {
		for _, location := range server.Location {
			// location = admin-app
			// httpUrl = http://127.0.0.1:8001/admin-app/v1/ap/
			if strings.Contains(allPath, location.LocationPath) {
				lp.ServerName = server.ServerName
				if location.LocationPath != "" && location.LocationPath != "/" {
					lp.ProxyPath = location.ProxyPath + strings.Split(allPath, location.LocationPath)[1]
				} else {
					lp.ProxyPath = strings.Split(allPath, server.ServerName)[1]
				}
			}
		}
		if lp.ProxyPath != "" {
			LoadProxyUrlMap[allPath] = lp
			return lp
		}
	}
	return LoadProxy{}
}
