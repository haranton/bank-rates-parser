package sender

import (
	"bank-rates-parser/internal/config"
	"bank-rates-parser/internal/pb"
	"bank-rates-parser/internal/scraper"
	"bank-rates-parser/internal/storage"
	"context"
	"fmt"
	"log/slog"
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
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			s.ParseAndSendMessage(ctx)
		}
	}()
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
}

func (s *Sender) Send(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
