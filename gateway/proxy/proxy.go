package proxy

import (
	"eff-gateway/discovery"
	"eff-gateway/gateway/proxy/types"
	"eff-gateway/glog"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	ProxyMap = make(map[string]types.Proxy)
)

var adminUrl = flag.String("adminUrl", "/api/v1/", "admin的地址")
var profile = flag.String("profile", "", "环境")
var proxyFile = flag.String("proxyFile", "/api/v1/", "测试环境的数据")

func InitProxy() {
	if *profile != "" {
		all, err := discovery.EC.Get("/")
		if err != nil {
			glog.ErrorLog.Fatalln("etcd service error:", err.Error())
		}
		for _, v := range all {
			t := types.Proxy{}
			err := json.Unmarshal(v.Value, &t)
			if err != nil {
				continue
			}
			ProxyMap[string(v.Key)] = t
		}
	} else {
		glog.InfoLog.Printf("加载本地配置数据: %s", *proxyFile)
		LoadProxyListFromFile()
	}

}

func InitProxyList() {
	resp, _ := http.Get(*adminUrl)
	if resp != nil && resp.StatusCode == 200 {
		bytes, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("ioutil.ReadAll err=", err)
			return
		}
		var response types.Response
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			fmt.Println("json.Unmarshal err=", err)
			return
		}
		proxyList := response.Data
		for _, proxy := range proxyList {
			//追加 反斜杠，为了动态匹配的时候 防止 /proxy/test  /proxy/test1 无法正确转发
			ProxyMap[proxy.Prefix+"/"] = proxy
		}
	}
}

// Forward
// 请求代理的实现函数
// @Param w 响应写入
// @Param r http请求
func Forward(ProxyMap map[string]types.Proxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			io.WriteString(w, "Request path Error")
			return
		}
		//从内存里面获取转发的url
		//去掉开头
		var value types.Proxy
		var ok bool
		for s, proxy := range ProxyMap {
			if strings.HasPrefix(r.RequestURI, s) {
				ok = true
				value = proxy
			}
		}
		if !ok {
			io.WriteString(w, "404 Gateway error")
			return
		}

		upstream := suffixURI(value)
		if value.RewritePrefix != "" {
			r.URL.Path = strings.ReplaceAll(r.RequestURI, value.Prefix, value.RewritePrefix)
		}
		// 解析url
		remote, err := url.Parse(value.IPAddr)
		glog.InfoLog.Printf("RequestURI:%s upstream:%s remote:%s", r.RequestURI, upstream, remote)
		if err != nil {
			glog.InfoLog.Println(err)
		}

		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Header.Set("Upgrade", r.Header.Get("websocket"))
		r.Header.Set("Connection", r.Header.Get("Upgrade"))

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(w, r)
	}

}
func LoadProxyListFromFile() {
	file, err := os.Open(*proxyFile)
	if err != nil {
		glog.ErrorLog.Println("err:", err)
	}
	var respond types.Response
	// 创建json解码器
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&respond)
	if err != nil {
		fmt.Println("LoadProxyListFromFile failed", err.Error())
	}
	proxyList := respond.Data
	for _, proxy := range proxyList {
		// 拼接的 key 例子：
		// proxy.Prefix = ？
		// proxy.Prefix+"/" = http://127.0.0.1:8087/admin
		ProxyMap[proxy.Prefix+"/"] = proxy
	}
}

// suffixURI 拆分重组请求地址
// requestURI http 的请求URI
// httpPath   http.URL.Path
func suffixURI(value types.Proxy) string {

	//从内存里面获取转发的url
	var upstream = ""
	//如果首位不是/开头，则需要追加
	if !strings.HasPrefix(value.RewritePrefix, "/") && value.RewritePrefix != "" {
		upstream += "/" + value.RewritePrefix
	} else {
		upstream += value.RewritePrefix
	}
	//如果转发的地址是 / 结尾的,需要去掉
	if strings.HasSuffix(value.Upstream, "/") {
		upstream += strings.TrimRight(value.Upstream, "/")
	} else {
		upstream += value.Upstream
	}

	return upstream
}
