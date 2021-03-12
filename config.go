package api

import "fmt"

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
}

func (c *Config) GetMysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Database)
}
