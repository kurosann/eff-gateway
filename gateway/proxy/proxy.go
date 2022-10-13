package proxy

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-gateway/glog"
	"go-gateway/src/proxy/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	proxyMap = make(map[string]types.Proxy)
)

var adminUrl = flag.String("adminUrl", "/api/v1/", "admin的地址")
var profile = flag.String("profile", "", "环境")
var proxyFile = flag.String("proxyFile", "/api/v1/", "测试环境的数据")

func init() {
	proxyMap["clean/add/"] = types.Proxy{Prefix: "", Upstream: "/api/v1/clean/add", RewritePrefix: ""}
	if *profile != "" {
		glog.InfoLog.Printf("加载远端数据: %s", *adminUrl)
		InitProxyList()
	} else {
		glog.InfoLog.Printf("加载本地配置数据: %s", *proxyFile)
		LoadProxyListFromFile()
	}
}

func Forward(writer http.ResponseWriter, request *http.Request) {
	HostReverseProxy(writer, request)
}

func InitProxyList() {
	resp, _ := http.Get(*adminUrl)
	if resp != nil && resp.StatusCode == 200 {
		bytes, err := ioutil.ReadAll(resp.Body)
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
			proxyMap[proxy.Prefix+"/"] = proxy
		}
	}
}

// HostReverseProxy
// 请求代理的实现函数
// w 响应写入
// r http请求
func HostReverseProxy(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		io.WriteString(w, "Request path Error")
		return
	}
	//从内存里面获取转发的url
	//去掉开头
	upstream := ""
	value, ok := proxyMap[r.RequestURI]
	if ok {
		upstream = suffixURI(value)
		r.URL.Path = strings.ReplaceAll(r.URL.Path, r.RequestURI, "")
	}
	// parse the url
	remote, err := url.Parse(upstream)
	glog.InfoLog.Printf("RequestURI %s upstream %s remote %s", r.RequestURI, upstream, remote)
	if err != nil {
		panic(err)
	}

	r.URL.Host = remote.Host
	r.URL.Scheme = remote.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = remote.Host

	httputil.NewSingleHostReverseProxy(remote).ServeHTTP(w, r)
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
		proxyMap[proxy.Prefix+"/"] = proxy
	}
}

// suffixURI 拆分重组请求地址
// requestURI http 的请求URI
// httpPath   http.URL.Path
func suffixURI(value types.Proxy) string {

	//从内存里面获取转发的url
	var upstream = ""
	//如果转发的地址是 / 结尾的,需要去掉
	if strings.HasSuffix(value.Upstream, "/") {
		upstream += strings.TrimRight(value.Upstream, "/")
	} else {
		upstream += value.Upstream
	}

	//如果首位不是/开头，则需要追加
	if !strings.HasPrefix(value.RewritePrefix, "/") {
		upstream += "/" + value.RewritePrefix
	} else {
		upstream += value.RewritePrefix
	}
	return upstream
}
