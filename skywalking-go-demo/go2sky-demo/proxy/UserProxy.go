package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	//_ "github.com/apache/skywalking-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http"
	"net/http/httputil"
	"net/url"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"time"
)

var (
	rpgrpc        bool
	rpOapServer   string
	rpOapAuth     string
	rpListenAddr  string
	rpServiceName string
	rpTargetAddr  string
	rpClient      *http.Client
	//db            *sql.DB
)

func init() {
	flag.BoolVar(&rpgrpc, "grpc", true, "use grpc reporter")
	//9.223.77.222:11800 需替换为 APM 的私网接入点
	flag.StringVar(&rpOapServer, "oap-server", "", "oap server address")
	flag.StringVar(&rpOapAuth, "oap-auth", "", "oap server auth")
	flag.StringVar(&rpListenAddr, "listen-addr", "localhost:10000", "listen address")
	flag.StringVar(&rpServiceName, "service-name", "test_longxi_proxy", "service name")
	flag.StringVar(&rpTargetAddr, "target-addr", "localhost:9999", "target service address")
	flag.Parse()
}

func test1() {
	r, err := reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()
	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	// This for test
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatalf("create new local span error %v \n", err)
	}
	span.SetOperationName("invoke data")
	span.Tag("kind", "outer")
	time.Sleep(time.Second)
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		log.Fatalf("create new sub local span error %v \n", err)
	}
	subSpan.SetOperationName("invoke inner")
	subSpan.Log(time.Now(), "inner", "this is right")
	time.Sleep(time.Second)
	subSpan.End()
	time.Sleep(500 * time.Millisecond)
	span.End()
	time.Sleep(time.Second)
}

func getOperationName(c *gin.Context) string {
	return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}

func proxy() {

	report, err := reporter.NewGRPCReporter(
		rpOapServer,
		reporter.WithAuthentication(rpOapAuth))
	if err != nil {
		log.Fatalf("crate grpc reporter error: %v \n", err)
	}
	defer report.Close()

	tracer, err := go2sky.NewTracer(rpServiceName, go2sky.WithReporter(report))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 注册中间件，将tracer实例注入到请求中
	r.Use(func(c *gin.Context) {
		span, ctx, err := tracer.CreateEntrySpan(c.Request.Context(), getOperationName(c), func(key string) (string, error) {
			return c.Request.Header.Get(key), nil
		})
		if err != nil {
			c.Next()
			return
		}

		span.SetComponent(5006)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(agentv3.SpanLayer_Http)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()
	})

	r.GET("/*path", func(c *gin.Context) {
		u := &url.URL{}
		//转发的协议，如果是https，写https.
		u.Scheme = "http"
		u.Host = rpTargetAddr
		proxy := httputil.NewSingleHostReverseProxy(u)

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("http: proxy error: %v", err)
		}

		// Inject context into HTTP request header `sw8`
		span, err := tracer.CreateExitSpan(c.Request.Context(), "/test/forwarding", rpTargetAddr, func(key, value string) error {
			c.Request.Header.Set(key, value)
			return nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		span.SetComponent(5006)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(agentv3.SpanLayer_Http)

		//c.Request = c.Request.WithContext(ctx)

		proxy.ServeHTTP(c.Writer, c.Request)

		//span.Tag(go2sky.TagStatusCode, strconv.Itoa(resp.StatusCode))
		span.End()

	})

	r.Run(rpListenAddr)
}

func main() {
	proxy()
}
