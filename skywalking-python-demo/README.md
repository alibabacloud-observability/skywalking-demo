# SkyWalking上报Python应用数据Demo
## 简介
1. 一个包含http请求和操作MySql的Python Demo。
2. 使用SkyWalking-Python Agent埋点监控，并将数据上报至可观测链路OpenTelemetry。
## 获取接入点信息
1. 登录[链路追踪控制台](https://tracing.console.aliyun.com/)，在页面顶部选择需要接入的区域。
2. 在左侧导航栏单击集群配置，然后单击接入点信息页签。
3. 在接入点信息页签的集群信息区域打开显示Token开关。
4. 在客户端采集工具区域单击SkyWalking，获取接入点信息。
## 用SkyWalking为Python应用自动埋点
### SkyWalking-Python版本 1.0.1
1. 下载package
```
pip3 install flask
pip3 install flask_sqlalchemy
pip3 install requests=2.19.1
pip3 install apache-skywalking==1.0.1
```
2. 在proxy/TestProxy.py和controller/TestController.py中分别配置接入信息和应用信息。
```
# 说明：<endpoint>和<auth-token需要登录控制台，进入可观测链路OpenTelemetry版控制台的接入流程中查询>
from skywalking import config
config.init(agent_collector_backend_services=<endpoint>,
            agent_authentication=<auth-token>,
            agent_name=<service-name>,
```
2. 加载skywalking-agent。
```
from skywalking import agent
agent.start()
```
3. 启动demo，并产生上报数据。
```
# demo用到了flask_sqlalchemy库操作MySql，所以需要配置数据库的连接信息SQLALCHEMY_DATABASE_URI
app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = ''
db = SQLAlchemy(app)

# demo采用Flask实现简单的请求转发，http请求先到达TestProxy，main再将请求转发给TestController。所以需要启动两个python文件。
python3.11 proxy/TestProxy.py --host <YourHostName1> --port <YourPort1>
python3.11 controller/TestController.py --host <YourHostName2> --port <YourPort2>
```
4. 使用Http请求，触发监控数据生成及上报。
```
# 方法1: 命令行使用crul构建http请求。
curl -l http://127.0.0.1:9999/forward/user

# 方法2: 在httprequests子文件夹中用requests库发送http请求。示例demo在APITest.py中。
import requests
url = 'http://127.0.0.1:9999/forward/user'
username = 'LiuHan'
password = 'lh232r425'
headers = {'Content-Type': 'application/json'}
data = json.dumps({'username': username, 'password': password})
response = requests.post(url, headers=headers, data=data)

# 方法3：在APITest.http中构建http请求。
GET http://127.0.0.1:9999/forward/user
```

