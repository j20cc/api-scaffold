package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"time"
)

var (
	MysqlCli *gorm.DB
	RedisCli *redis.Client
	err      error
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func InitDB() {
	initMysql()
	initRedis()
}

func initRedis() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pwd"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := RedisCli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

var dsn = "%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local"

func initMysql() {
	MysqlCli, err = gorm.Open("mysql", fmt.Sprintf(dsn,
		viper.GetString("mysql.user"),
		viper.GetString("mysql.pwd"),
		viper.GetString("mysql.addr"),
		viper.GetString("mysql.db")))
	if err != nil {
		log.Fatal(err)
	}
	//defer DB.Close()

	MysqlCli.AutoMigrate(&User{})
}
