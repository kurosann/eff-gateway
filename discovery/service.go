package discovery

import (
	"eff-gateway/clients"
	"eff-gateway/setting"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var EC *clients.EtcdClient

func InitService() error {
	endpoint := setting.Config.Etcd.Endpoint
	options := clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Duration(setting.Config.Etcd.DialTimeout) * time.Millisecond,
	}
	EC = clients.NewEtcdClient(options)
	return nil
}

func KeepAlive(key string, f clients.EtcdEvFunc) {
	EC.KeepWatch(key, f)
}

func DelSv(key string) error {
	return EC.Del(key)
}
func PutSv(key string, value string) error {
	return EC.Put(key, value)
}
