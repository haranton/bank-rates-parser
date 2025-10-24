package main

import (
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.GetLogger(cfg.Env)

	application := app.NewApp(logger, cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go application.GRPCSrv.MustRun()

	<-stop

	// // Создаем экземпляр скрапера
	// scraper := deposit_scraper.NewScraper()

	// // Инициализируем WebDriver
	// err := scraper.Initialize()
	// if err != nil {
	// 	log.Fatalf("Ошибка инициализации: %v", err)
	// }
	// defer scraper.Close()

	// // Собираем данные о вкладах
	// cards, err := scraper.ScrapeDeposits()
	// if err != nil {
	// 	log.Fatalf("Ошибка сбора данных: %v", err)
	// }

	// // Выводим результаты
	// scraper.PrintCards(cards)

	// fmt.Println("Проверка завершена успешно!")
}
