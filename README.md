# GoFlux
## 项目目录树

```shell
├── logs
└── src
    ├── config
    ├── internal
    │   ├── core
    │   │   ├── grpc
    │   │   └── http
    │   └── example
    │       ├── gin_server
    │       │   ├── beans
    │       │   ├── cmd
    │       │   ├── config
    │       │   ├── controller
    │       │   ├── middleware
    │       │   ├── routes
    │       │   └── serializer
    │       └── http_client
    │           └── simple_demo
    └── pkg
        ├── e
        │   └── constant
        └── utils
            ├── dbutil
            └── logutil
```


## V1.x 版本概述

go的快速通信框架 flux旨在形容 数据在框架中流动自如 

http客户端和grpc连接池客户端：

1）go语言基于的net/http封装的快速开发万能框架

2）grpc连接池封装与请求 优化连接池

3）提供样例example,给出基于gin框架的服务端，并提供客户端，给出使用案例

待补充...