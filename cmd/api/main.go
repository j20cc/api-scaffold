package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lukedever/api"
	"github.com/lukedever/api/http"
	"github.com/lukedever/api/mysql"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	cfgFile = flag.String("config", "config.yaml", "config file path")

	config api.Config
)

func main() {
	flag.Parse()

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}

	// new db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.User, config.Mysql.Password, config.Mysql.Addr, config.Mysql.Database)
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// new userservice
	userService := mysql.NewUserService(db)

	srv := http.NewServer(&config)
	srv.UserService = userService

	go func() {
		if err := srv.Run(); err != nil && errors.Is(err, netHttp.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
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
