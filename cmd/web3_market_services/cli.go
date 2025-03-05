package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/version"
	"github.com/urfave/cli/v2"

	"github.com/qiaopengjun5162/web3-market-services/common/cliapp"
	"github.com/qiaopengjun5162/web3-market-services/common/opio"
	"github.com/qiaopengjun5162/web3-market-services/config"
	"github.com/qiaopengjun5162/web3-market-services/database"
	flags2 "github.com/qiaopengjun5162/web3-market-services/flags"
	"github.com/qiaopengjun5162/web3-market-services/services"
)

// Semantic holds the textual version string for major.minor.patch.
var Semantic = fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

// WithMeta holds the textual version string including the metadata.
var WithMeta = func() string {
	v := Semantic
	if version.Meta != "" {
		v += "-" + version.Meta
	}
	return v
}()

func withCommit(gitCommit, gitDate string) string {
	vsn := WithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	if (version.Meta != "stable") && (gitDate != "") {
		vsn += "-" + gitDate
	}
	return vsn
}

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	fmt.Println("running rpc server...")
	cfg := config.NewConfig(ctx)

	grpcServerCfg := &services.MarketRpcConfig{
		Host: cfg.RpcServer.Host,
		Port: cfg.RpcServer.Port,
	}

	db, err := database.NewDB(context.Background(), cfg.MasterDB)
	if err != nil {
		log.Error("Failed to create db", "err", err)
		return nil, err
	}
	return services.NewMarketRpcService(grpcServerCfg, db)
}

func runMigrations(ctx *cli.Context) error {
	ctx.Context = opio.CancelOnInterrupt(ctx.Context)
	log.Info("Running migrations...")
	cfg := config.NewConfig(ctx)
	db, err := database.NewDB(ctx.Context, cfg.MasterDB)
	if err != nil {
		log.Error("Failed to connect to database", "err", err)
		return err
	}
	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Error("Failed to close database", "err", err)
		}
	}(db)
	return db.ExecuteSQLMigration(cfg.Migrations)
}

func NewCli(GitCommit string, gitDate string) *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Name:                 "Web3 market services",
		Usage:                "Web3 market services",
		Description:          "An market services with rpc",
		Version:              withCommit(GitCommit, gitDate),
		EnableBashCompletion: true, // Boolean to enable bash completion commands
		Commands: []*cli.Command{
			{
				Name:        "rpc",
				Usage:       "Start rpc server",
				Description: "Start rpc server",
				Flags:       flags,
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "migrate",
				Usage:       "Migrate database",
				Description: "Migrate database",
				Flags:       flags,
				Action:      runMigrations,
			},
			{
				Name:        "version",
				Usage:       "Show project version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
