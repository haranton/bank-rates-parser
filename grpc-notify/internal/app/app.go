package app

import (
	grpcApp "grpc-notify/internal/app/grpc"
	"grpc-notify/internal/config"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcApp.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	config *config.Config,
) *App {
	grpcApp := grpcApp.NewApp(log, config)

	return &App{
		GRPCSrv: grpcApp,
	}
}
