整体架构思路

etcd 启动

```shell
etcd --name=node1 \
 --listen-client-urls http://127.0.0.1:2379 \
 --advertise-client-urls http://127.0.0.1:2379 \
 --listen-peer-urls http://127.0.0.1:2380 \
 --initial-advertise-peer-urls http://127.0.0.1:2380 \
 --initial-cluster node1=http://127.0.0.1:2380,node2=http://127.0.0.1:2480,node3=http://127.0.0.1:2580 \
 --initial-cluster-state new --initial-cluster-token token-1
 
etcd --name=node2 \
 --listen-client-urls http://127.0.0.1:2479 \
 --advertise-client-urls http://127.0.0.1:2479 \
 --listen-peer-urls http://127.0.0.1:2480 \
 --initial-advertise-peer-urls http://127.0.0.1:2480 \
 --initial-cluster node1=http://127.0.0.1:2380,node2=http://127.0.0.1:2480,node3=http://127.0.0.1:2580 \
 --initial-cluster-state new --initial-cluster-token token-1

etcd --name=node3 \
 --listen-client-urls http://127.0.0.1:2579 \
 --advertise-client-urls http://127.0.0.1:2579 \
 --listen-peer-urls http://127.0.0.1:2580 \
 --initial-advertise-peer-urls http://127.0.0.1:2580 \
 --initial-cluster node1=http://127.0.0.1:2380,node2=http://127.0.0.1:2480,node3=http://127.0.0.1:2580 \
 --initial-cluster-state new --initial-cluster-token token-1
 
```

添加一个节点
```shell

etcdctl --endpoints=127.0.0.1:2379,127.0.0.1:2479,127.0.0.1:2579 member add node4 --peer-urls=http://127.0.0.1:2180

etcd --name node4 --listen-client-urls http://127.0.0.1:2179 \
--advertise-client-urls http://127.0.0.1:2179 \
--listen-peer-urls http://127.0.0.1:2180 \
--initial-advertise-peer-urls http://127.0.0.1:2180 \
--initial-cluster-state existing \
--initial-cluster node1=http://127.0.0.1:2380,node2=http://127.0.0.1:2480,node3=http://127.0.0.1:2580,node4=http://127.0.0.1:2180 \
--initial-cluster-token token-1
 
```

![](docs/imgs/架构图.png)