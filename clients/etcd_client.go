package clients

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
)

type EtcdEvFunc func(event *clientv3.Event) bool

type EtcdClient struct {
	client *clientv3.Client

	locker  sync.RWMutex
	isClose bool
}

func NewEtcdClient(options clientv3.Config) *EtcdClient {
	client, err := clientv3.New(options)
	if err != nil {
		log.Println("etcd run failed:", err.Error())
	}
	return &EtcdClient{
		client:  client,
		locker:  sync.RWMutex{},
		isClose: false,
	}
}

func (c *EtcdClient) Close() {
	c.locker.Lock()
	defer c.locker.Unlock()

	if !c.isClose {
		c.isClose = true
		_ = c.client.Close()
	}
}

// Put 插入值
func (c *EtcdClient) Put(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := c.client.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}

// Del 删除值
func (c *EtcdClient) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := c.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// GetFirst 获取一个值
func (c *EtcdClient) GetFirst(key string) (*mvccpb.KeyValue, error) {
	all, err := c.Get(key)
	if len(all) == 0 {
		return nil, err
	}
	return all[0], err
}

// Get 获取
func (c *EtcdClient) Get(key string) ([]*mvccpb.KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	response, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return response.Kvs, nil
}

// KeepWatch 订阅
func (c *EtcdClient) KeepWatch(key string, f EtcdEvFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	watch := c.client.Watch(ctx, key)
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
