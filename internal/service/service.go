package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	api "github.com/dagherghinescu/companies/internal/http"
	"github.com/dagherghinescu/companies/internal/http/routes"
	"github.com/dagherghinescu/companies/internal/logger"
)

// Service holds the application dependencies and configuration.
type Service struct {
	Log    *zap.Logger
	APICfg *api.Config
}

// New creates a new Service instance, initializing logger and configuration.
// Returns an error if the context is canceled or configuration fails.
func New(ctx context.Context) (*Service, error) {
	if ctx.Err() != nil {
		return nil, errors.New("context canceled")
	}

	configs, err := validateConfigs()
	if err != nil {
		return nil, err
	}

	logger, err := logger.Init()
	if err != nil {
		panic(err)
	}

	return &Service{
		Log:    logger,
		APICfg: configs.httpSrv,
	}, nil
}

func Run(ctx context.Context, svc *Service) error {
	r := gin.Default()
	routes.RegisterCompanyRoutes(r)

	srv := &http.Server{
		Addr:              svc.APICfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: svc.APICfg.ReadHeaderTimeout,
		ReadTimeout:       svc.APICfg.ReadTimeout,
		WriteTimeout:      svc.APICfg.WriteTimeout,
	}

	go func() {
		if err := api.StartServer(ctx, svc.Log, srv); err != nil && err != http.ErrServerClosed {
			svc.Log.Error("HTTP server stopped with error", zap.Error(err))
		} else {
			svc.Log.Info("HTTP server stopped")
		}
	}()

	svc.Log.Info("Application is running", zap.String("addr", svc.APICfg.Addr))
	return nil
}

// Close releases resources held by Service
func (d *Service) Close() {
	if d.Log != nil {
		// Ensure all buffered logs are written
		_ = d.Log.Sync()
	}
}
