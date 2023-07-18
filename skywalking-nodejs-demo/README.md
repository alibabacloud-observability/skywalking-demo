# SkyWalking上报Nodejs应用数据Demo
## 简介
1. 一个包含http请求代理转发和操作MySql的Go Demo。
2. 使用SkyWalking Nodejs Agent埋点监控，并将数据上报至可观测链路OpenTelemetry。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为Nodejs应用自动埋点
1. 下载package
```
npm install express body-parser
npm install --save skywalking-backend-js
```
2. 运行Demo
```
# 说明：<endpoint>和<token>对应接入点信息。

# step1: 配置环境变量。
export SW_AGENT_COLLECTOR_BACKEND_SERVICES=<endpoint>
export SW_AGENT_AUTHENTICATION=<token>
export SW_AGENT_NAME=<sevice-name>

# step2: 进入demo根目录，启动demo。
node index.js
```
3. 使用Http请求，触发监控数据生成及上报。
```
# 方法1: 命令行使用curl
curl -X GET http://localhost:3000/api/service1
curl -X GET http://localhost:3000/api/service2

# 方法2: 打开浏览器，在地址栏输入url
http://localhost:3000/api/service1
http://localhost:3000/api/service2

# 方法3: Intellij IDEA 新建http文件，在里面写http请求请求
```
4. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，查看应用监控状态。