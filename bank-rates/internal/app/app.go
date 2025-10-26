package app

import (
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/db"
	"bank-rates-parser/internal/scraper"
	"bank-rates-parser/internal/sender"
	"bank-rates-parser/internal/storage"
	"context"
	"log/slog"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	sender *sender.Sender
	scrp   *scraper.Scraper
}

func New(cfg *config.Config, slogger *slog.Logger) *App {

	dbConn := db.GetDBConnect(cfg, slogger)
	db.MustRunMigrations(cfg, slogger)

	storage := storage.NewStorage(dbConn)
	scrp := scraper.NewScraper(slogger)

	send, err := sender.NewSender(storage, scrp, slogger, cfg)
	if err != nil {
		panic(err)
	}

	return &App{
		cfg:    cfg,
		logger: slogger,
		sender: send,
		scrp:   scrp,
	}
}

func (app *App) Start(ctx context.Context) {
	app.sender.Start(ctx)
}

func (app *App) Close() {
	app.scrp.Close()
}
