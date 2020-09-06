package main

import (
	"flag"
	"os"
	"strings"

	"gvue-scaffold/app/models"
	"gvue-scaffold/cmd"
	"gvue-scaffold/cmd/migrate"
	"gvue-scaffold/internal/log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/urfave/cli/v2"
)

var (
	cfg = "etc/config.yml"
)

func main() {
	flag.Parse()
	//初始化配置
	initConfig()
	//初始化日志
	log.NewLogger()
	//初始化db
	models.InitDB()

	app := &cli.App{
		Name:  viper.GetString("app.name"),
		Usage: "gin+vue全栈开发",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run web server",
				Action: func(c *cli.Context) error {
					cmd.RunServer()
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					return migrate.RunUp(c)
				},
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "create migrations, run 'migrate create [create|alter]_foos_table'",
						Action: func(c *cli.Context) error {
							return migrate.RunCreate(c)
						},
					},
					{
						Name:  "rollback",
						Usage: "rollback migrations",
						Action: func(c *cli.Context) error {
							return migrate.RunUp(c)
						},
					},
					{
						Name:  "refresh",
						Usage: "refresh migrations",
						Action: func(c *cli.Context) error {
							return migrate.RunUp(c)
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("run app error", zap.Error(err))
	}
}

func initConfig() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		panic(".env file is not exists")
	}
	_ = godotenv.Load()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GVUE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(cfg)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//TODO:检查必要的配置
	if viper.GetString("app.locale") == "" {
		viper.SetDefault("app.locale", "zh")
	}
}
