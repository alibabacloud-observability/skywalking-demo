# SkyWalking上报Rust应用数据Demo
## 简介
1. 一个简单的Rust Demo，基于Hyper构建的Http服务器。
2. 使用skywalking-rust agent埋点监控，并将数据上报至可观测链路OpenTelemetry。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为Rust应用埋点
1. Rust项目引入skywalking-rust
```
# 前提条件：安装protobuf
## Debian-base系统
sudo apt install protobuf-compiler
## MacOS系统
brew install protobuf

# 项目添加skywalking-rust
cargo add skywalking --features vendored

# 引入头文件
use skywalking::{reporter::grpc::GrpcReporter, trace::tracer::Tracer};
```
2. 手动埋点
```
# 说明: skywalking-rust暂不支持监控应用自动埋点，所以需要使用者手动埋点。

# 使用EntrySpan、LocalSpan、ExitSpan完成调用链手动埋点，使用这三种Span，可以完成全链路追踪。
EntrySpan: 入口Span，服务端使用EntrySpan从Http请求中拿到链路追踪上下文。
LocalSpan: 本地Span，同一进程内可使用该Span埋点。
ExitSpan: 出口Span，客户端使用ExitSpan向Http请求中注入链路追踪上下文。
```
3. 设置接入点信息
```
# 说明: <EndPoint>是接入点地址，<TOKEN>是接入点token，<ServiceName>是应用名称。

let collector_backend_address = "<EndPoint>";
let auth_token = "<TOKEN>";
let service_name = "<ServiceName>";

let reporter = GrpcReporter::connect(collector_backend_address).await?;
let reporter = reporter.with_authentication(auth_token);
let tracer = Tracer::new(service_name, "instance", reporter.clone());
```
4. 使用Http请求，触发监控数据生成及上报。
```
# 方法1: 命令行使用curl
curl -X GET http://localhost:9999/hello
curl -X POST -H "Content-Type: application/json" -d '{"name": "Alice", "age": 30}' http://localhost:9999/echo

# 方法2: 打开浏览器，在地址栏输入url
http://localhost:9999/hello

# 方法3: idea新建http文件，在里面写http请求
```
5. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，查看应用监控状态。