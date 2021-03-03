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
	"github.com/lukedever/api/http"
	"github.com/lukedever/api/mysql"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mode    = flag.String("mode", "debug", "app mode, debug test or release")
	addr    = flag.String("addr", "127.0.0.1:9001", "listen addr")
	cfgFile = flag.String("config", "config.yaml", "config file path")

	config Config
)

func main() {
	flag.Parse()

	readConfig()

	// new db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.User, config.Mysql.Password, config.Mysql.Addr, config.Mysql.Database)
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// new userservice
	userService := mysql.NewUserService(db)

	srv := http.NewServer(*mode, *addr)
	srv.UserService = userService

	go func() {
		if err := srv.Run(); err != nil && errors.Is(err, netHttp.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server exiting")
}

// Config is main config
type Config struct {
	Title string

	Mysql struct {
		Addr     string
		User     string
		Password string
		Database string
	}
}

func readConfig() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}
}
