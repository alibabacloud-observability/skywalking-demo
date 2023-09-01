// require('make-promises-safe')

const {default: agent} = require("skywalking-backend-js");
// import agent from 'skywalking-backend-js';

// require('make-promises-safe')

// agent.start({
//     // serviceName: "NodeJs-Agnet",
//     // serviceInstance: "NodeJs-Agent",
//     // collectorAddress: "10.114.126.48:11800"
// });
agent.start({
});



const express = require('express');
const bodyParser = require('body-parser');
const app = express();

// 解析请求体
app.use(bodyParser.json());

// 定义路由映射
const routes = {
    '/api/service1': 'service1',
    '/api/service2': 'service2',
    // 可以添加更多路由映射
};

// 处理请求
app.all('*', (req, res) => {
    const service = routes[req.path];
    if (service) {
        // 将请求转发给service处理
        const serviceModule = require(`./services/${service}`);
        serviceModule.handle(req, res);
    } else {
        res.status(404).send('Not Found');
    }
});

// 启动服务
app.listen(3000, () => {
    console.log('Server is running on port 3000');
});