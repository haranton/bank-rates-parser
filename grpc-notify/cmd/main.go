package main

import (
	"grpc-notify/internal/app"
	"grpc-notify/internal/config"
	"grpc-notify/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.GetLogger(cfg.Env)

	application := app.NewApp(logger, cfg.GRPC.Port, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go application.GRPCSrv.MustRun()

	<-stop

	application.GRPCSrv.Stop()
	logger.Info("app succesfully stop")

}
