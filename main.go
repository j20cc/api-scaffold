package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gvue-scaffold/app/controllers"
	"gvue-scaffold/app/middlewares"
	"gvue-scaffold/app/models"
	"gvue-scaffold/pkg/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	cfg          = flag.String("f", "etc/config.yml", "specified the config file name")
	staticFolder = "./public"
	indexFile    = "./public/index.html"
)

func main() {
	flag.Parse()
	//初始化配置
	initConfig(*cfg)
	//初始化日志
	log.NewLogger()
	//初始化db
	models.InitDB()

	runMode := viper.GetString("app.mode")
	runAddr := viper.GetString("app.addr")
	gin.SetMode(runMode)
	router := gin.Default()
	//注册路由
	registerRoutes(router)
	srv := &http.Server{
		Addr:           runAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("err occurred at app starting...", zap.Error(err))
			panic(err)
		}
	}()
	log.Info("app is running...", zap.String("addr", runAddr))

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown err:", zap.Error(err))
	}
	log.Error("server exited")
}

func registerRoutes(r *gin.Engine) {
	//前端
	r.Use(static.Serve("/", static.LocalFile(staticFolder, true)))
	r.LoadHTMLGlob(indexFile)
	//跨域
	r.Use(middlewares.CORS())
	//注册路由
	baseController := new(controllers.Controller)
	r.GET("/api/hello", baseController.Hello)
	//auth-route-start
	userController := new(controllers.User)
	r.POST("/api/register", userController.Register)
	r.POST("/api/login", userController.Login)
	r.POST("/api/password/email", userController.SendResetEmail)
	r.POST("/api/password/reset", userController.ResetPassword)
	auth := r.Group("/api", middlewares.Auth())
	auth.POST("/verification/email", userController.SendVerifyEmail)
	auth.POST("/verification", userController.VerifyEmail)
	auth.GET("/profile", userController.GetProfile)
	//auth-route-end
}

func initConfig(file string) {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		panic(".env file is not exists")
	}
	_ = godotenv.Load()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GVUE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(file)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//TODO:检查必要的配置
	if viper.GetString("app.locale") == "" {
		viper.SetDefault("app.locale", "zh")
	}
}
