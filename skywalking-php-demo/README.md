# SkyWalking上报PHP应用数据Demo
## 简介
1. 一个简单的PHP Demo，借助swoole实现http server。
2. 使用SkyWalking PHP Agent埋点监控，并将数据上报至可观测链路OpenTelemetry。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为PHP应用自动埋点
1. 安装skywalking-php agent
```
前提条件：
    PHP 7.2 - 8.x.
    GCC
    Rustc 1.56+
    Cargo
    Libclang 9.0+
    Make
    Protoc

安装skywalking-php agent：
    pecl install skywalking_agent
    
    或者（从源代码编译）
    
    git clone --recursive https://github.com/apache/skywalking-php.git
    cd skywalking-php
    
    phpize
    ./configure
    make
    make install

验证是否启动成功
    php --ri skywalking_agent
    
    预期输出：
    skywalking_agent

    version => 0.5.0
    authors => Apache Software Foundation:jmjoy <jmjoy@apache.org>:Yanlong He <heyanlong@apache.org>
    
```

* 补充：demo使用了swoole框架，mac安装swoole请参考[安装Swoole](https://www.easyswoole.com/QuickStart/installSwoole.html)

3. 配置skywalking-php agent参数（将配置参数添加到你的php.ini文件中）

* your-service-name: 你的服务名称。
* your-backend-server-address: 你的接入点地址。
* your-auth-token: 你的接入点鉴权令牌。

```
[skywalking_agent]
extension = skywalking_agent.so

; Enable skywalking_agent extension or not.
skywalking_agent.enable = 1

; Log file path.
skywalking_agent.log_file = /tmp/skywalking-agent.log

; Log level: one of `OFF`, `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`.
skywalking_agent.log_level = INFO

; Address of skywalking oap server.
skywalking_agent.server_addr = <your-backend-server-address>
skywalking_agent.authentication = <your-auth-token>

; Application service name.
skywalking_agent.service_name = <your-service-name>

; 说明：demo使用了swoole框架，所以需要额外安装swoole依赖
[swoole]
extension=swoole.so
```
3. 重启PHP应用
4. 使用Http请求，触发监控数据生成及上报，
```
浏览器地址栏输入 http://localhost:10000/ 或 http://localhost:10000/ping
```
5. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，查看应用监控状态。