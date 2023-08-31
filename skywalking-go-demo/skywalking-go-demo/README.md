# 基于skywalking-go agent上报Go应用数据
## 简介
1. 一个包含http_server和http_client的demo
2. 使用新版探针skywalking-go监控go demo
3. client文件夹，包含http_client
4. server文件夹，包含http_server
5. config文件夹，包含agent_config.yaml Agent属性配置文件
6. skywalking-go-0.2.0文件夹，为skywalking-go探针。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 接入流程
1. go下载agent依赖。
```
go get -u github.com/apache/skywalking-go
```
2. 配置Agent属性相关的环境变量。说明：也可在config/agent_config.yaml中配置。
```
# 说明：<service_name>表示服务名称，<endpoint>表示接入点，<token>表示接入点鉴权令牌。后两者按照”获取接入点信息“流程得到。
export SW_AGENT_NAME=<service_name>
export SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=<endpoint>
export SW_AGENT_REPORTER_GRPC_AUTHENTICATION=<token>
```
3. 构建项目
```
# server
cd server && sudo go build -toolexec "/path/to/skywalking-go-0.2.0/bin/skywalking-go-agent--<ostype>-amd64 -config /path/to/config/agent_config.yaml" -a

# client
cd client && go build
```
4. 启动应用。说明：先启动server后启动client。
5. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，查看应用监控状态。