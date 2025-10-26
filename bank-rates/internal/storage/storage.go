package storage

import (
	"bank-rates-parser/internal/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// CreateBankRate добавляет новую запись о банке.
func (st *Storage) CreateUpdateBankRate(ctx context.Context, bankRate models.BankRate) error {
	query := `
		INSERT INTO bank_rates (bank_name, deposit_name, rate)
		VALUES ($1, $2, $3)
		ON CONFLICT (bank_name) DO UPDATE
		SET deposit_name = EXCLUDED.deposit_name,
		    rate = EXCLUDED.rate
	`

	_, err := st.db.ExecContext(ctx, query,
		bankRate.BankName,
		bankRate.DepositName,
		bankRate.Rate,
	)
	if err != nil {
		return fmt.Errorf("failed to create or update bank rate: %w", err)
	}

	return nil
}

// BankRates возвращает все записи из таблицы bank_rates.
func (st *Storage) BankRates(ctx context.Context) ([]models.BankRate, error) {
	var rates []models.BankRate

	query := `
		SELECT id, bank_name, deposit_name, rate
		FROM bank_rates
		ORDER BY id
	`

	if err := st.db.SelectContext(ctx, &rates, query); err != nil {
		return nil, fmt.Errorf("failed to get bank rates: %w", err)
	}

	return rates, nil
}
