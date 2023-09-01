# 链路追踪SkyWalking上报数据demo

## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。


## 用SkyWalking为Java应用自动埋点

###  SkyWalking版本 8.x
1. 打开agent-8.x/config/agent.config，配置接入信息和应用信息。
```
collector.backend_service=<endpoint> 
agent.authentication=<auth-token>
agent.service_name=<ServiceName>
```
例如使用SkyWalking 8.\*的版本，同时选择需要接入的区域为杭州，则配置
```
collector.backend_service=tracing-analysis-dc-hz.aliyuncs.com:8000
```

2. 运行命令启动样例程序和SkyWalking探针

```
java -javaagent:<work-path>/agent-8.x/skywalking-agent.jar -jar <work-path>/demo.jar
```
\<work-path\>是当前示例程序所在路径。

3. 构造上报数据，运行下面命令，产生测试数据

```
curl -L localhost:8080/hello
```
4. 在控制台检查应用有没有创建成功以及数据有没有上报成功

## 补充
```
根目录下skywalking-agent/为skywalking-java v8.16.0版本包
```
