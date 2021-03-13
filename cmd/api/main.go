package main

import (
	"fmt"

	"github.com/lukedever/api"
	"github.com/lukedever/api/http"
	"github.com/lukedever/api/mysql"
	"github.com/spf13/viper"
)

var (
	config api.Config
)

func main() {
	readConfig()

	// new db
	db, err := mysql.NewDB(config.GetMysqlDsn(), config.Mode)
	if err != nil {
		panic(err)
	}

	// new userservice
	userService := mysql.NewUserService(db)

	// new server
	svr := http.NewServer(&config)
	svr.UserService = userService

	svr.Run()
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
