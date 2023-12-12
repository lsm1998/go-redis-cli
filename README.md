# go_redis_cli

redis client implemented by golang

## 自己动手造轮子系列（一）

使用golang语言实现的redis客户端工具

多年前写的一个小工具， 核心代码不到100行，实现了Redis Resp协议的解析，支持基本的交互命令和发布订阅功能

其他系列

. [自己动手造轮子系列（二）- 基于Netty实现tomcat](https://github.com/lsm1998/tiny_tomcat)

### 环境说明

>= go 1.15

## 使用方法

### 1. 编译
````shell
go build -o redis_cli
````

### 2. 运行
````shell
./redis_cli -h 127.0.0.1 -p 6379 -pass youpassword
````

## 功能

### 1. 基本交互命令

![基本命令.png](/doc/基本命令.png)

![基本命令.png](/doc/认证.png)

### 2. 发布订阅

![发布订阅.png](/doc/发布订阅.png)
