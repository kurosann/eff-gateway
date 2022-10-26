整体架构思路

etcd 启动
```shell
etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380
```



![](docs/imgs/架构图.png)