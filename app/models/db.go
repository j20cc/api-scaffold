package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var (
	mysqlCli *gorm.DB
	redisCli *redis.Client
	err      error
)

// Model is base model
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// InitDB init mysql and redis
func InitDB() {
	initMysql()
	initRedis()
}

func initRedis() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pwd"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

var dsn = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"

func initMysql() {
	mysqlCli, err = gorm.Open("mysql", fmt.Sprintf(dsn,
		viper.GetString("mysql.user"),
		viper.GetString("mysql.pwd"),
		viper.GetString("mysql.addr"),
		viper.GetString("mysql.db")))
	if err != nil {
		log.Fatal(err)
	}
	//defer DB.Close()

	mysqlCli.AutoMigrate(&User{})
}
