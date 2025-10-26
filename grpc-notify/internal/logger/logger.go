package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// GetLogger возвращает цветной логгер в зависимости от окружения
func GetLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "DEBUG":
		handler := tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "15:04:05",
			NoColor:    false, // включаем цвета
		})
		logger = slog.New(handler)

	case "PRODUCTION":
		handler := tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    false,
		})
		logger = slog.New(handler)

	default:
		// fallback — обычный текстовый хэндлер
		handler := tint.NewHandler(os.Stdout, &tint.Options{
			Level:   slog.LevelInfo,
			NoColor: false,
		})
		logger = slog.New(handler)
	}

	return logger
}
