# SkyWalking上报.NET应用数据Demo
## 简介
1. 一个通过SkyWalking .NET Core SDK上报.NET应用数据的Demo
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为.NET应用自动埋点
1. 安装.NET运行环境，参考[.NET官网](https://dotnet.microsoft.com/zh-cn/download/dotnet/thank-you/sdk-7.0.306-macos-arm64-installer)，SkyWalking .NET Core SDK支持netcoreapp3.1、net5.0、net6.0以及更高版本，该Demo使用.NET 7.0 SDK (v7.0.306)
2. 添加SkyWalking .NET Core SDK，根目录下执行
```
dotnet add package SkyAPM.Agent.AspNetCore
```
3. 在根目录的skyapm.json中配置skywalking参数。
```
<your-service-name>: 你的服务/应用名称。
<your-collector-backend-address>: 你的接入点地址，参考获取接入点信息。
<your-collect-auth-token>: 你的接入点鉴权令牌，参考获取接入点信息。
```
4. 配置环境变量。
```
export ASPNETCORE_HOSTINGSTARTUPASSEMBLIES=SkyAPM.Agent.AspNetCore
export SKYWALKING__SERVICENAME=<your-service-name>
```
5. 启动应用。
```
dotnet run
```
6. 访问 http://localhost:5276 ，多刷新几次，生成上报数据，登录OpenTelemetry控制台，选择具体区域，查看目标应用。
