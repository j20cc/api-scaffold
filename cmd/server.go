package cmd

import (
	"context"
	"gvue-scaffold/app/controllers"
	"gvue-scaffold/app/middlewares"
	"gvue-scaffold/internal/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	staticFolder = "public"
	indexFile    = "public/index.html"
)

// RunServer run web server
func RunServer() {
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
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("err occurred at app starting...", zap.Error(err))
		}
	}()
	log.Info("app is running...", zap.String("addr", runAddr))

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown error", zap.Error(err))
	}
	log.Info("server exited")
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
