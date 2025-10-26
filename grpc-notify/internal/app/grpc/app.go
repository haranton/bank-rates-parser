package grpcApp

import (
	"fmt"
	"grpc-notify/internal/config"
	"grpc-notify/internal/grpc/notify"
	"grpc-notify/internal/service"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	config     *config.Config
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	gRPCServer := grpc.NewServer()
	srv := service.NewTelegramClient(config.Token, config.ChatId, log)
	notify.Register(gRPCServer, srv)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		config:     config,
	}
}

func (a *App) Run() error {
	const op = "grpcApp"

	log := a.log.With(slog.String("op", op))

	log.Info("start gRPC server", "port", a.config.GRPC.Port)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// serve gRPC
	if err := a.gRPCServer.Serve(l); err != nil {
		log.Error("gRPC server stopped with error", "err", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server stopped")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.config.GRPC.Port))

	a.gRPCServer.GracefulStop()

}
