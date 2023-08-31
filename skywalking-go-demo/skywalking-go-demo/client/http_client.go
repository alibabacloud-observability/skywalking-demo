package main

import (
	"fmt"
	_ "github.com/apache/skywalking-go"
	"net/http"
	"time"
)

func main() {
	// 设置间隔时间
	interval := 50 * time.Millisecond // 每隔 50 毫秒发送一次请求

	// 设置目标服务器地址
	url := "http://localhost:9999"

	// 创建一个无限循环，定时发送请求
	for {
		// 发送 HTTP GET 请求
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("请求发送失败:", err)
		} else {
			fmt.Println("请求已发送，状态码:", resp.StatusCode)
			resp.Body.Close()
		}

		// 等待指定的间隔时间
		time.Sleep(interval)
	}
}
