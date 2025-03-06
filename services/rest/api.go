package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/qiaopengjun5162/web3-market-services/common/httputil"
	"github.com/qiaopengjun5162/web3-market-services/config"
	"github.com/qiaopengjun5162/web3-market-services/database"
	"github.com/qiaopengjun5162/web3-market-services/services/rest/routes"
	"github.com/qiaopengjun5162/web3-market-services/services/rest/service"
)

const (
	HealthCheckPath  = "/health"
	SupportAssetPath = "/api/v1/get_support_asset"
	MarketPricePath  = "/api/v1/get_market_price"
)

type ApiConfig struct {
	HttpServerConfig    config.ServerConfig
	MetricsServerConfig config.ServerConfig
}

type Api struct {
	router    *chi.Mux
	apiServer *httputil.HTTPServer
	db        *database.DB
	stopped   atomic.Bool
}

func NewApi(ctx context.Context, cfg *config.Config) (*Api, error) {
	api := &Api{}
	if err := api.initFromConfig(ctx, cfg); err != nil {
		return nil, errors.Join(err, api.Stop(ctx))
	}
	return api, nil
}

func (a *Api) initFromConfig(ctx context.Context, cfg *config.Config) error {
	if err := a.initDB(ctx, cfg); err != nil {
		log.Error("Failed to init db", "err", err)
		return fmt.Errorf("failed to init db: %w", err)
	}
	a.initRouter(cfg.RestServer, cfg)
	if err := a.startServer(cfg.RestServer); err != nil {
		log.Error("Failed to start server", "err", err)
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (a *Api) initDB(ctx context.Context, cfg *config.Config) error {
	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("Failed to init db", "err", err)
		return fmt.Errorf("failed to init db: %w", err)
	}
	a.db = db
	return nil
}

func (a *Api) initRouter(_ config.ServerConfig, _ *config.Config) {
	v := new(service.Validator)
	svc := service.NewHandleSvc(v, a.db.MarketPrice, a.db.OfficialCoinRate)
	apiRouter := chi.NewRouter()
	r := routes.NewRoutes(apiRouter, svc)

	apiRouter.Use(middleware.Timeout(time.Second * 12))
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(middleware.RequestID)
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Heartbeat(HealthCheckPath))

	apiRouter.Get(SupportAssetPath, r.GetSupportAsset)
	apiRouter.Get(MarketPricePath, r.GetMarketPrice)

	a.router = apiRouter
}

func (a *Api) Start(ctx context.Context) error {
	return nil
}

func (a *Api) Stop(ctx context.Context) error {
	var result error
	if a.apiServer != nil {
		if err := a.apiServer.Stop(ctx); err != nil {
			result = errors.Join(fmt.Errorf("failed to stop api server: %w", err), result)
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			result = errors.Join(fmt.Errorf("failed to close db: %w", err), result)
		}
	}
	a.stopped.Store(true)
	log.Info("API stopped")
	return result
}

func (a *Api) startServer(serverConfig config.ServerConfig) error {
	log.Debug("API server listening...", "port", serverConfig.Port)
	addr := net.JoinHostPort(serverConfig.Host, strconv.Itoa(serverConfig.Port))
	log.Info("API server listening...", "addr", addr)
	srv, err := httputil.NewHTTPServer(addr, a.router)
	if err != nil {
		return fmt.Errorf("failed to start API server: %w", err)
	}
	log.Info("API server started", "addr", srv.Addr().String())
	a.apiServer = srv
	return nil
}

func (a *Api) Stopped() bool {
	return a.stopped.Load()
}
