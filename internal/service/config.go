package service

import (
	"fmt"

	api "github.com/dagherghinescu/companies/internal/http"
	"github.com/dagherghinescu/companies/internal/http/middleware"
	"github.com/dagherghinescu/companies/internal/kafka"
	"github.com/dagherghinescu/companies/internal/repository"
)

type config struct {
	httpSrv  *api.Config
	dbCfg    *repository.Config
	jwtCfg   *middleware.JWTConfig
	kafkaCfg *kafka.Config
}

func validateConfigs() (*config, error) {
	srvConfig, err := api.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("server configuration error: %w", err)
	}

	pgCfg, err := repository.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("database configuration error: %w", err)
	}

	jwtCfg, err := middleware.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("jwt secret error: %w", err)
	}

	kafkaCfg, err := kafka.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("kafka config error: %w", err)
	}

	return &config{
		httpSrv:  srvConfig,
		dbCfg:    pgCfg,
		jwtCfg:   jwtCfg,
		kafkaCfg: kafkaCfg,
	}, nil
}
