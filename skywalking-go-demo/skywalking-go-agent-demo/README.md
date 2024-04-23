# SkyWalking Go Agent上报 Go 应用数据 Demo

1. 下载 [skywalking go agent](https://www.apache.org/dyn/closer.cgi/skywalking/go/0.4.0/apache-skywalking-go-0.4.0-bin.tgz)
* 解压后的 skywalking agent 所在目录为 `/path/to/apache-skywalking-go-0.4.0-bin/bin`，不同操作系统的 agent 名称不同：
  * macOS: skywalking-go-agent-0.4.0-darwin-amd64
  * Linux: skywalking-go-agent-0.4.0-linux-amd64
  * Windows: skywalking-go-agent-0.4.0-windows-amd64

2. 下载 skywalking-go 所需依赖项，在 demo 根目录下执行以下命令：
```
go get github.com/apache/skywalking-go
```

3. 在 main package 中导入 skywalking module
> 注意：demo 中已经导入 skywalking module，这一步可跳过

* 方法一：手动添加
```go
package main

import (
	_ "github.com/apache/skywalking-go"
)
```

* 方法二: 通过 skywalking-agent 注入
```
/path/to/skywalking-go-agent -inject path/to/your-project
```
  

4. 使用以下命令构建项目
```
go build -toolexec="/path/to/skywalking-go-agent" -a -o main .
```

5. 运行项目
```
./main
```

6. 访问 Go 应用
```
curl http://localhost:8080/hello
```