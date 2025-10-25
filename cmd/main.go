package main

import (
	"bank-rates-parser/internal/app"
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.GetLogger(cfg.Env)

	application := app.New(cfg, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.TODO()

	go application.Start(ctx)

	<-stop
	application.Close()
	logger.Info("app succesfully stop")

}
