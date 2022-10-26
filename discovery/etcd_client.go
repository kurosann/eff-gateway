package discovery

import (
	"context"
	"eff-gateway/glog"
	"eff-gateway/setting"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
	"time"
)

var EC EtcdClient

type EtcdEvFunc func(event *clientv3.Event) bool

func InitEtcd() {
	endpoint := setting.Config.Etcd.Endpoint
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Duration(setting.Config.Etcd.DialTimeout) * time.Millisecond,
	})
	glog.InfoLog.Printf("services exited:%v\n", endpoint)
	if err != nil {
		glog.ErrorLog.Println(err.Error())
		glog.ErrorLog.Fatalln("etcd service is run fail!")
		return
	}
	glog.InfoLog.Println("[INFO] connect to etcd success")
	EC = EtcdClient{
		client:  c,
		locker:  sync.RWMutex{},
		isClose: false,
	}
}

type EtcdClient struct {
	client *clientv3.Client

	locker  sync.RWMutex
	isClose bool
}

func (c *EtcdClient) Close() {
	c.locker.Lock()
	defer c.locker.Unlock()

	if !c.isClose {
		c.isClose = true
		c.client.Close()
	}
}

func (c *EtcdClient) Put(key string, value string) error {
	_, err := c.client.Put(context.Background(), key, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *EtcdClient) Del(key string) error {
	_, err := c.client.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}
func (c *EtcdClient) GetFirst(key string) (*mvccpb.KeyValue, error) {
	all, err := c.Get(key)
	if len(all) == 0 {
		return nil, err
	}
	return all[0], err
}

func (c *EtcdClient) Get(key string) ([]*mvccpb.KeyValue, error) {
	response, err := c.client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	return response.Kvs, nil
}

func (c *EtcdClient) KeepWatch(key string, f EtcdEvFunc) {
	watch := c.client.Watch(context.Background(), key)
	for response := range watch {
		for _, ev := range response.Events {
			if f(ev) {
				continue
			} else {
				break
			}
		}
	}
}
