package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/lukedever/api/http"
	"github.com/lukedever/api/mysql"
	"gopkg.in/yaml.v2"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mode    = flag.String("mode", "dev", "app mode, dev test or production")
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

	svr := http.NewServer(*mode, *addr)
	svr.UserService = userService

	_ = svr.Run()
}

// Config is main config
type Config struct {
	Addr string

	Mysql struct {
		Addr     string
		User     string
		Password string
		Database string
	}
}

func readConfig() {
	file, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
}
