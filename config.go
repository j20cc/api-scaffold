package api

import (
	"fmt"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// Config is main config
type Config struct {
	Title  string
	Mode   string
	Addr   string
	Locale string

	Mysql struct {
		Host     string `mapstructure:"mysql_host"`
		Port     int    `mapstructure:"mysql_port"`
		User     string `mapstructure:"mysql_user"`
		Password string `mapstructure:"mysql_password"`
		Database string `mapstructure:"mysql_database"`
	} `mapstructure:",squash"`

	JWT struct {
		Secret string `mapstructure:"jwt_secret"`
		TTL    int    `mapstructure:"jwt_ttl"`
	} `mapstructure:",squash"`
}

func (c *Config) GetMysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database)
}
