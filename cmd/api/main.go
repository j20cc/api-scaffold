package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukedever/api"
	"github.com/lukedever/api/http"
	"github.com/lukedever/api/mysql"
	"github.com/spf13/viper"

	netHttp "net/http"
)

var (
	config api.Config
)

func main() {
	readConfig()
	log.Println(config.Mysql.Password)

	// new db
	db, err := mysql.NewDB(config.GetMysqlDsn(), config.Mode)
	if err != nil {
		panic(err)
	}

	// new userservice
	userService := mysql.NewUserService(db)

	// new server
	srv := http.NewServer(&config)
	srv.UserService = userService

	runAndwatchSignal(srv)
}

func readConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
}

func runAndwatchSignal(srv *http.Server) {
	go func() {
		if err := srv.Run(); err != nil && errors.Is(err, netHttp.ErrServerClosed) {
			log.Printf("listen error: %s\n", err)
		}
	}()
	log.Printf("server running on %s", config.Addr)

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exiting")
}
