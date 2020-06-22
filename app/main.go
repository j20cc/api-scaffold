package main

import (
	"flag"
	"github.com/lukedever/gvue-scaffold/app/controllers"
	"github.com/lukedever/gvue-scaffold/app/middlewares"
	"github.com/lukedever/gvue-scaffold/app/models"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	cfg          = flag.String("f", "config.yml", "specified the config file name")
	staticFolder = "./public"
	indexFile    = "./public/index.html"
)

func main() {
	flag.Parse()
	//初始化配置
	initConfig(*cfg)
	//初始化db
	models.InitDB()

	runMode := viper.GetString("app.mode")
	runAddr := viper.GetString("app.addr")
	gin.SetMode(runMode)
	r := gin.Default()
	registerRoutes(r)
	_ = r.Run(runAddr)
}

func registerRoutes(r *gin.Engine) {
	//前端
	r.Use(static.Serve("/", static.LocalFile(staticFolder, true)))
	r.LoadHTMLGlob(indexFile)
	//认证路由
	r.Use(middlewares.CORS())
	userController := new(controllers.User)
	r.POST("/api/register", userController.Register)
	r.POST("/api/login", userController.Login)
	r.POST("/api/password/email", userController.SendResetEmail)
	r.POST("/api/password/reset", userController.ResetPassword)
	auth := r.Group("/api", middlewares.Auth())
	auth.POST("/verification/email", userController.SendVerifyEmail)
	auth.POST("/verification", userController.VerifyEmail)
	//获取详细信息
	auth.GET("/profile", userController.GetProfile)
}

func initConfig(file string) {
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
}
