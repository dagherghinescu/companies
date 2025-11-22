package service

import (
	"fmt"

	api "github.com/dagherghinescu/companies/internal/http"
	"github.com/dagherghinescu/companies/internal/repository"
)

type config struct {
	httpSrv *api.Config
	dbCfg   *repository.Config
}

func validateConfigs() (*config, error) {
	srvConfig, err := api.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("server configuration error: %w", err)
	}

	pgCfg, err := repository.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("server configuration error: %w", err)
	}

	return &config{
		httpSrv: srvConfig,
		dbCfg:   pgCfg,
	}, nil
}
