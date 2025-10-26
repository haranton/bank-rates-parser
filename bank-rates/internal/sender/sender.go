package sender

import (
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/pb"
	"bank-rates-parser/internal/scraper"
	"bank-rates-parser/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Sender struct {
	storage *storage.Storage
	scraper *scraper.Scraper
	logger  *slog.Logger
	config  *config.Config
	client  pb.NotificationServiceClient
}

func NewSender(storage *storage.Storage, scraper *scraper.Scraper, logger *slog.Logger, cfg *config.Config) (*Sender, error) {

	conn, err := grpc.NewClient(cfg.NotifyServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to create grpc client", "error", err)
		return nil, err
	}
	client := pb.NewNotificationServiceClient(conn)

	return &Sender{
		storage: storage,
		scraper: scraper,
		logger:  logger,
		config:  cfg,
		client:  client,
	}, nil
}

func (s *Sender) Start(ctx context.Context) {
	tickerScraping := time.NewTicker(15 * time.Second)

	go func() {
		for range tickerScraping.C {
			s.ParseAndSendMessage(ctx)
		}
	}()

	tickerAnalitic := time.NewTicker(24 * time.Hour)

	go func() {
		for range tickerAnalitic.C {
			s.SendAnalytics(ctx)
		}
	}()

}

func (s *Sender) SendAnalytics(ctx context.Context) {
	const op = "sender.SendAnalytics"

	existingRates, err := s.storage.BankRates(ctx)
	if err != nil {
		s.logger.Error("failed to get existing rates from db", "error", err, "op", op)
		return
	}

	// Проверяем, есть ли данные для отправки
	if len(existingRates) == 0 {
		s.logger.Info("no bank rates found for analytics", "op", op)
		return
	}

	// Форматируем сообщение в более читаемом виде
	msg := " **Аналитика ставок по вкладам**\n\n"
	msg += "| Банк | Вклад | Ставка |\n"
	msg += "|------|-------|--------|\n"

	for _, bankRate := range existingRates {
		// Используем правильное форматирование для float и экранируем специальные символы
		rowStr := fmt.Sprintf("| %s | %s | %.2f%% |\n",
			escapeMarkdown(bankRate.BankName),
			escapeMarkdown(bankRate.DepositName),
			bankRate.Rate)
		msg += rowStr
	}

	// Добавляем статистику
	msg += fmt.Sprintf("\n**Всего банков:** %d", len(existingRates))

	if err := s.Send(msg); err != nil {
		s.logger.Error("failed to send analytics", "error", err, "op", op)
	}
}

// Вспомогательная функция для экранирования Markdown-символов
func escapeMarkdown(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}

func (s *Sender) ParseAndSendMessage(ctx context.Context) {
	const op = "sender.ParseAndSendMessage"

	// 1. Получаем свежие ставки с сайта
	bankRates, err := s.scraper.ScrapeDeposits()
	if err != nil {
		s.logger.Error("failed to scrape rates", "error", err, "op", op)
		return
	}

	// 2. Получаем существующие ставки из БД
	existingRates, err := s.storage.BankRates(ctx)
	if err != nil {
		s.logger.Error("failed to get existing rates from db", "error", err, "op", op)
		return
	}

	// 3. Преобразуем существующие ставки в map для удобства сравнения
	existingMap := make(map[string]float32)
	for _, r := range existingRates {
		existingMap[r.BankName] = r.Rate
	}

	// 4. Проверяем изменения и реагируем
	for _, newRate := range bankRates {
		oldRate, exists := existingMap[newRate.BankName]

		if !exists {
			// Новый банк — просто добавляем
			if err := s.storage.CreateUpdateBankRate(ctx, newRate); err != nil {
				s.logger.Error("failed to insert new bank rate", "bank", newRate.BankName, "error", err)
				continue
			}

			msg := fmt.Sprintf("Банк добавлен: %s со ставкой %.2f%%", newRate.BankName, newRate.Rate)
			if err := s.Send(msg); err != nil {
				s.logger.Error("failed to send new bank notification", "error", err)
			}
			continue
		}

		// Если ставка изменилась
		if oldRate != newRate.Rate {
			if err := s.storage.CreateUpdateBankRate(ctx, newRate); err != nil {
				s.logger.Error("failed to update bank rate", "bank", newRate.BankName, "error", err)
				continue
			}

			msg := fmt.Sprintf("Изменение ставки в банке %s: было %.2f%% → стало %.2f%%",
				newRate.BankName, oldRate, newRate.Rate)
			if err := s.Send(msg); err != nil {
				s.logger.Error("failed to send rate change notification", "error", err)
			}
		}
	}
	s.scraper.Close()
}

func (s *Sender) Send(msg string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := s.client.SendNotification(ctx, &pb.NotificationRequest{
		Type: "telegram",
		// Target:  s.config.Telegram.ChatID,
		Target:  "",
		Message: msg,
	})
	if err != nil {
		return fmt.Errorf("rpc error: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("send failed: %s", resp.Error)
	}

	return nil
}
