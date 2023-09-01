package main

import (
	"flag"
	"fmt"
	"github.com/SkyAPM/go2sky"
	pHttp "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
	sHttp "net/http"
	_ "net/http/httptest"
	"os"
	"time"
)

var (
	grpc      bool
	oapServer string
	oapAuth   string
	//ListenAddr  string
	serviceName string
	targetAddr  string
	client      *sHttp.Client
	interval    int
	//db            *sql.DB
)

func init() {
	flag.BoolVar(&grpc, "grpc", true, "use grpc reporter")
	//9.223.77.222:11800 需替换为 APM 的私网接入点
	flag.StringVar(&oapServer, "oap-server", "", "oap server address")
	flag.StringVar(&oapAuth, "oap-auth", "", "oap server auth")
	flag.StringVar(&serviceName, "service-name", "test_longxi_http_client", "service name")
	flag.StringVar(&targetAddr, "target-addr", "localhost:10000", "target service address")
	flag.IntVar(&interval, "req-interval", 1, "generate http request interval")

	flag.Parse()
}

func test1() error {
	report, err := reporter.NewGRPCReporter(
		oapServer,
		reporter.WithAuthentication(oapAuth))
	if err != nil {
		log.Fatalf("crate grpc reporter error: %v \n", err)
	}
	defer report.Close()

	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(report))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}
	client, err := pHttp.NewClient(tracer)

	for {
		url := fmt.Sprintf("http://%s/execute", targetAddr)

		req, err := sHttp.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("unable to create http request: %+v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("unable to do http request: %+v", err)
		}

		defer resp.Body.Close()

		// 将响应体直接打印到控制台
		if _, err := fmt.Fprintln(os.Stdout, resp.Body); err != nil {
			log.Fatalf("unable to print response body: %+v", err)
		}
		// 延迟固定间隔
		time.Sleep(time.Duration(interval) * time.Second)
	}

	return nil
}

func main() {
	test1()
}
