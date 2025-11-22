package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/dagherghinescu/companies/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svc, err := service.New(ctx)
	if err != nil {
		log.Printf("could not initialize service: %+v", err)
		return
	}
	defer svc.Close()

	if err := service.Run(ctx, svc); err != nil {
		log.Printf("could not start service: %+v", err)
		return
	}

	svc.Log.Info("Application is running")
	<-ctx.Done()
	svc.Log.Info("Shutting down application")
}
