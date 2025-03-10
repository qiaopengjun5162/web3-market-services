package main

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"

	"github.com/qiaopengjun5162/web3-market-services/common/opio"
)

var (
	GitCommit = ""
	gitDate   = ""
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	app := NewCli(GitCommit, gitDate)
	ctx := opio.WithInterruptBlocker(context.Background())
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("Application failed")
		os.Exit(1)
	}
}
