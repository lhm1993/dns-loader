# dns-loader
go语言实现的dns负载测试工具

### 1 启动程序

#### 1.1 参数说明

```
Usage: dns-loader [options...] 
Options:
  -t       loader type, one of "master","worker","once"
  -c       config file path for app start
  -s       dns server
  -p       dns server listen port. Default is 53.
  -d       query domain name
  -D       duration time. Default 60 seconds
  -q       query per second. Default is 10
  -r       random subdomain length. Default is 5
  -R       enable random query type. Default is false
  -Q       query type. Default is A
  -debug   enable debug mode
```


#### 1.1  单次执行模式
 
该模式下无需指定配置文件和发包模式，所有的配置通过命令行来输入，仅仅执行一次发包，单机模式，如需要中断直接运行ctrl+C即可退出。

```
go run dns-loader.go -s 172.17.0.2 -d jsmean.com -Q A -D 60 -q 1000 -r 5

2017/11/15 20:58:04 new dns loader client success, start send packet...
2017/11/15 20:58:04 Start generate dns packet
2017/11/15 20:58:04 init dns loader client configuration success
...
[===>                ] 
...
2017/11/15 20:59:04 [Result]total packets sum:60001
2017/11/15 20:59:04 [Result]runing time 1m0s
2017/11/15 20:59:04 [Result]status nxdomain:60001 [100.00%]


```
#### 1.2  master模式

该模式下需要指定模式类型为`master`和配置文件，运行后通过web浏览器来管理发包请求，默认登入用户名和密码为admin/admin,登入配置发包类型和模式后点击开始即可，该模式下可以添加agent，只需要agent端运行agent模式即可（见1.3）
```
go run dns-loader.go -t master -c config.ini 
2017/11/15 21:11:46 load configuration from file:config.ini
2017/11/15 21:11:46 start Web for control panel default web address:localhost:9889
2017/11/15 21:11:46 http server route init success
2017/11/15 21:11:46 static file folder:/web/assets
...
```
启动后访问：http://localhost:9889

#### 1.3  agent模式

该模式需要指定模式类型为`agent`和配置文件，可放入后台运行，该模式不会开启任何的web页面，但是会提供http接口供master服务器端调用
```
go run dns-loader.go -t agent -c config.ini 
2017/11/15 21:15:12 load configuration from file:config.ini
2017/11/15 21:15:12 start agent server listen on 8998 for master connect
2017/11/15 21:15:12 agent server route init success
...
```

### 2 性能优化

##### 函数性能

使用原生库来展示调用栈

```
go tool pprof -seconds=5 localhost:8080/debug/pprof/profile
> web
```

使用uber的库来展示调用栈
```
go get github.com/uber/go-torch
git clone https://github.com/brendangregg/FlameGraph.git
go-torch -t 5
firefox torch.svg

```

