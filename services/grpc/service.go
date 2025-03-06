package grpc

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"

	"github.com/qiaopengjun5162/web3-market-services/database"
	"github.com/qiaopengjun5162/web3-market-services/protobuf/market"
)

const MAX_RECEIVE_MESSAGE_SIZE = 1024 * 1024 * 30000

type MarketRpcConfig struct {
	Host string
	Port int
}

type MarketRpcServices struct {
	*MarketRpcConfig

	db *database.DB

	market.UnimplementedMarketServicesServer
	stopped atomic.Bool
}

func NewMarketRpcService(cfg *MarketRpcConfig, db *database.DB) (*MarketRpcServices, error) {
	return &MarketRpcServices{
		MarketRpcConfig: cfg,
		db:              db,
	}, nil
}

func (ms *MarketRpcServices) Start(ctx context.Context) error {
	go func(ms *MarketRpcServices) {
		rpcAddr := fmt.Sprintf("%s:%d", ms.MarketRpcConfig.Host, ms.MarketRpcConfig.Port)
		listener, err := net.Listen("tcp", rpcAddr)
		if err != nil {
			log.Error("Failed to listen rpc server", "err", err)
			return
		}

		opt := grpc.MaxRecvMsgSize(MAX_RECEIVE_MESSAGE_SIZE)
		grpcServer := grpc.NewServer(opt, grpc.ChainUnaryInterceptor(nil))
		reflection.Register(grpcServer)
		market.RegisterMarketServicesServer(grpcServer, ms)

		log.Info("Starting rpc server", "addr", listener.Addr())
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("Failed to start rpc server", "err", err)
		}
	}(ms)
	return nil
}

func (ms *MarketRpcServices) Stop(ctx context.Context) error {
	ms.stopped.Store(true)
	return nil
}

func (ms *MarketRpcServices) Stopped() bool {
	return ms.stopped.Load()
}
