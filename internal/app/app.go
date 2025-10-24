package app

import (
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/db"
	"bank-rates-parser/internal/scraper"
	"bank-rates-parser/internal/storage"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	cfg     *config.Config
	logger  *slog.Logger
	storage *storage.Storage
	scraper *scraper.Scraper
}

func New(cfg *config.Config, slogger *slog.Logger) (*App, error) {
	dbConn := db.GetDBConnect(cfg, slogger)
	if err := db.RunMigrations(cfg, slogger); err != nil {
		return nil, err
	}

	storage := storage.NewStorage(dbConn)
	f := scraper.NewScraper()

	svc := service.NewService(rep, f)
	s := sender.NewSender(svc, b, slogger)
	hand := handler.NewHandler(svc, b, slogger)

	return &App{
		cfg:        cfg,
		logger:     slogger,
		httpClient: &http.Client{Timeout: 35 * time.Second},
	}, nil
}
