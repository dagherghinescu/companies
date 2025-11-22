package service

import (
	"fmt"

	api "github.com/dagherghinescu/companies/internal/http"
)

type config struct {
	httpSrv *api.Config
}

func validateConfigs() (*config, error) {
	srvConfig, err := api.EnvConfig()
	if err != nil {
		return nil, fmt.Errorf("server configuration error: %w", err)
	}

	return &config{
		httpSrv: srvConfig,
	}, nil
}
