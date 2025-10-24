package db

import (
	"bank-rates-parser/internal/config"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
)

func GetDBConnect(config *config.Config, logger *slog.Logger) *sqlx.DB {

	var DBHost string
	if config.Env == "PRODUCTION" {
		DBHost = config.Database.HostProd
	} else {
		DBHost = config.Database.Host
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBHost,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	logger.Info("Connecting to database", slog.String("dsn", dsn))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect database")
		os.Exit(1)
	}

	return db

}
