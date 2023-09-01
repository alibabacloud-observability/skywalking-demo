# SkyWalking上报Go应用数据Demo
## 简介
1. 一个包含http请求和操作MySql的Python Demo。
2. 使用SkyWalking-Python Agent埋点监控，并将数据上报至可观测链路OpenTelemetry。
3. service基于Gin框架实现，解析Http请求并对MySql数据库执行增删查改操作。
4. proxy接收Http请求并对其进行路由转发，其中实现了部分手动埋点。
5. client基于Go2Sky（Go2Sky实现了Http的自动埋点）实现Http客户端实例，周期向service发送Http请求。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为Go应用自动埋点
1. 下载package
```
# 下载skywalking go sdk：go2sky
go get -u github.com/SkyAPM/go2sky

# 下载go2sky的gin中间件
go get github.com/SkyAPM/go2sky-plugins/gin/v3

# 下载go2sky的sql中间件
go get github.com/SkyAPM/go2sky-plugins/sql
```
2. 启动demo产生上报数据
```
# 说明：service的监听地址是proxy的目标转发地址。上报数据相关的三个参数是oap-server、oap-auth、service-name。
# oap-server、oap-auth、service属于go2sky相关环境变量，也可以直接设置环境变量，go2sky会自动获取。

# step1: 启动service/UserService.go
cd skywalking-go-demo/service
go build
./service --dsn <your-sql-dsn> --listen-addr <your-bind-address> --oap-server <your-collector-backend-address> --oap-auth <your-oap-auth> --service-name <your-service-name> --sql-peer-addr <your-mysql-peer-addr>

# step2: 启动proxy/UserProxy.go
cd skywalking-go-demo/proxy
go build
./proxy --dsn <your-sql-dsn> --listen-addr <your-bind-address> --oap-server <your-collector-backend-address> --oap-auth <your-oap-auth> --service-name <your-service-name> --target-addr <your-target-forwarding-addr>

# 补充: 在client/HttpClient.go中实现了Http客户端（Go2Sky增强版）测试，Go2Sky支持对Http Client的自动埋点。启动方式:
cd skywalking-go-demo/client
go build
./client --oap-server <your-collector-backend-address> --oap-auth <your-oap-auth> --service-name <your-service-name> --target-addr <your-target-forwarding-addr>
```
3. 使用Http请求，触发监控数据生成及上报。
```
# 方法1: 命令行使用curl
curl -X GET http://<your-proxy-listen-addr>/execute

# 方法2: 打开浏览器，在地址栏输入url
http://<your-proxy-listen-addr>/execute

# 方法3: goland 新建http文件，在里面写http请求

# 补充:如果使用client/HttpClient，无需自行生成Http请求。
```
4. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，查看应用监控状态。