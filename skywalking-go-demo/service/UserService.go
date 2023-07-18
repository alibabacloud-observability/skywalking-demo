package main

import (
	"context"
	"flag"
	"fmt"
	_ "fmt"
	"github.com/SkyAPM/go2sky"
	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	sqlPlugin "github.com/SkyAPM/go2sky-plugins/sql"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jmoiron/sqlx"
	"log"
	"net/http"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"time"
)

var (
	rsgrpc        bool
	rsOapServer   string
	rsOapAuth     string
	rsListenAddr  string
	rsServiceName string
	dsn           string
	sqlPeerAddr   string
	rsClient      *http.Client
	//db            *sql.DB
)

type testFunc func(context.Context, *sqlPlugin.DB) error

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func getOperationName(c *gin.Context) string {
	return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}

func init() {
	flag.BoolVar(&rsgrpc, "grpc", true, "use grpc reporter")
	//9.223.77.222:11800 需替换为 APM 的私网接入点
	flag.StringVar(&rsOapServer, "oap-server", "", "oap server address")
	flag.StringVar(&rsOapAuth, "oap-auth", "", "oap server auth")
	flag.StringVar(&rsListenAddr, "listen-addr", "localhost:10000", "listen address")
	flag.StringVar(&rsServiceName, "service-name", "test_longxi_service", "service name")
	flag.StringVar(&dsn, "dsn", "root@tcp(127.0.0.1:3306)/skywalking_test", "mysql dsn")
	flag.StringVar(&sqlPeerAddr, "sql-peer-addr", "127.0.0.1:3306", "mysql peer addr")

	flag.Parse()
}

func db_test() {
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//var report go2sky.Reporter
	report, err := reporter.NewGRPCReporter(
		rsOapServer,
		reporter.WithAuthentication(rsOapAuth))

	if err != nil {
		log.Fatalf("crate grpc reporter error: %v \n", err)
	}

	defer report.Close()

	tracer, err := go2sky.NewTracer(rsServiceName, go2sky.WithReporter(report))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}

	db, err1 := sqlPlugin.Open("mysql", dsn, tracer,
		sqlPlugin.WithSQLDBType(sqlPlugin.MYSQL),
		sqlPlugin.WithQueryReport(),
		sqlPlugin.WithParamReport(),
		sqlPlugin.WithPeerAddr(sqlPeerAddr),
	)

	if err1 != nil {
		log.Fatalf("open db error: %v \n", err)
	}

	r.Use(v3.Middleware(r, tracer))

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

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "service_test")
	})

	r.GET("/execute", func(c *gin.Context) {
		tests := []struct {
			name string
			fn   testFunc
		}{
			{"exec", testExec},
			{"stmt", testStmt},
			{"commitTx", testCommitTx},
			{"rollbackTx", testRollbackTx},
		}
		for _, test := range tests {
			log.Printf("excute test case %s", test.name)
			if err1 := test.fn(c.Request.Context(), db); err1 != nil {
				log.Fatalf("test case %s failed: %v", test.name, err1)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "execute sql success",
		})

	})

	// 创建用户
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		result, err := db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", user.Username, user.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		user.Id = int(id)

		c.JSON(201, user)
	})

	// 查询单个用户
	r.GET("/users/:id", func(c *gin.Context) {
		var user User
		err := db.QueryRow("SELECT id, username, password FROM user WHERE id = ?", c.Param("id")).Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	})

	// 查询用户列表
	r.GET("/users", func(c *gin.Context) {
		var users []User
		rows, err := db.Query("SELECT id, username, password FROM user")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.Id, &user.Username, &user.Password)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(200, users)
	})

	// 更新用户
	r.PUT("/users/:id", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE user SET username = ?, password = ? WHERE id = ?", user.Username, user.Password, c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	})

	// 删除用户
	r.DELETE("/users/:id", func(c *gin.Context) {
		_, err := db.Exec("DELETE FROM user WHERE id = ?", c.Param("id"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Status(204)
	})

	r.Run(rsListenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func testExec(ctx context.Context, db *sqlPlugin.DB) error {
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS users`); err != nil {
		return fmt.Errorf("exec drop error: %w", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY, username VARCHAR(255), password VARCHAR(255))`); err != nil {
		return fmt.Errorf("exec create error: %w", err)
	}
	// test insert
	if _, err := db.ExecContext(ctx, `INSERT INTO users (username, password) VALUE (?, ?)`, "foo", "mypassword"); err != nil {
		return fmt.Errorf("exec insert error: %w", err)
	}
	var username string
	// test select
	if err := db.QueryRowContext(ctx, `SELECT username FROM users WHERE id = ?`, 1).Scan(&username); err != nil {
		return fmt.Errorf("query select error: %w", err)
	}
	fmt.Printf("Username: %s\n", username)
	return nil
}

func testStmt(ctx context.Context, db *sqlPlugin.DB) error {
	stmt, err := db.PrepareContext(ctx, `INSERT INTO users (username, password) VALUE (?, ?)`)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	if _, err := stmt.ExecContext(ctx, "bar", "mypassword2"); err != nil {
		return err
	}
	return nil
}

func testCommitTx(ctx context.Context, db *sqlPlugin.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx error: %v", err)
	}
	if _, err := tx.Exec(`INSERT INTO users (username, password) VALUE (?, ?)`, "foobar", "mypassword3"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET username = ? WHERE id = ?`, "foobar2", 1); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func testRollbackTx(ctx context.Context, db *sqlPlugin.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx error: %v", err)
	}
	if _, err := tx.Exec(`UPDATE users SET password = ? WHERE id = ?`, "mypassword4", 3); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET username = ? WHERE id = ?`, "foobar3", 2); err != nil {
		return err
	}
	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}
