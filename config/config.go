package config

import (
	"github.com/urfave/cli/v2"

	"github.com/qiaopengjun5162/web3-market-services/flags"
)

type Config struct {
	Migrations string
	RpcServer  ServerConfig
	MasterDB   DBConfig
	SlaveDB    DBConfig
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"dbname"`
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		Migrations: ctx.String(flags.MigrationsFlag.Name),
		RpcServer: ServerConfig{
			Host: ctx.String(flags.RpcHostFlag.Name),
			Port: ctx.Int(flags.RpcPortFlag.Name),
		},
		MasterDB: DBConfig{
			Host:     ctx.String(flags.MasterDbHostFlag.Name),
			Port:     ctx.Int(flags.MasterDbPortFlag.Name),
			User:     ctx.String(flags.MasterDbUserFlag.Name),
			Password: ctx.String(flags.MasterDbPasswordFlag.Name),
			Name:     ctx.String(flags.MasterDbNameFlag.Name),
		},
		SlaveDB: DBConfig{
			Host:     ctx.String(flags.SlaveDbHostFlag.Name),
			Port:     ctx.Int(flags.SlaveDbPortFlag.Name),
			User:     ctx.String(flags.SlaveDbUserFlag.Name),
			Password: ctx.String(flags.SlaveDbPasswordFlag.Name),
			Name:     ctx.String(flags.SlaveDbNameFlag.Name),
		},
	}
}
