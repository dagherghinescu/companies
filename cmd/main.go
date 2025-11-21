package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	api "github.com/dagherghinescu/companies/internal/http"
	"github.com/dagherghinescu/companies/internal/http/routes"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routes.RegisterCompanyRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := api.StartServer(ctx, logger, srv); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}
