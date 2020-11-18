package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"cm-v5/serv/urls"

	"github.com/dlintw/goconf"
	raven "github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
)

var (
	port string
)

// @title Swagger API Version 4
// @version 4.0
// @description ott-backend-v3 server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email thu.nguyen.thuy@dzones.vn

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1993
// @BasePath /backend/cm/v4

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	// Init Router
	router := initRouterDefault()

	// Url Router
	urls.InitUrlsV3(router)
	urls.InitUrlsWarmUp(router)
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	// router.Run()
	// For a hard coded port
	// router.Run(":" + port)

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Listening and serving HTTP on :", port)
	s.ListenAndServe()

}

func init() {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered: ", r)
	// 	}
	// }()

	flagvar_port := flag.String("port", "1993", "Set Port")
	flag.Parse()
	port = *flagvar_port
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Show info service
	fmt.Println("CONFIG: ")
	fmt.Println("- Port: " + port)
	fmt.Println("- Num CPU: " + fmt.Sprint(runtime.NumCPU()))

	initSentry()
}

func initSentry() {
	CommonConfig, err := goconf.ReadConfigFile("config/common.config")
	if err != nil {
		log.Println("ReadConfigFile", err)
	}
	SENTRY_DSN, _ := CommonConfig.GetString("SENTRY", "dsn")
	raven.SetDSN(SENTRY_DSN)
}

func initRouterDefault() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(sentry.Recovery(raven.DefaultClient, false))

	f, _ := os.Create("log/error.log")
	gin.DefaultErrorWriter = io.MultiWriter(f)
	return router
}

func initRouterDefaultWithLog() *gin.Engine {
	// Disable Console Color
	gin.DisableConsoleColor()

	// Log to file
	f, _ := os.Create("log/access.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.Use(sentry.Recovery(raven.DefaultClient, false))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	return router
}
