package models

type BankRate struct {
	ID          int64  `db:"id"`
	BankName    string `db:"bank_name"`
	DepositName string `db:"deposit_name"`
	Rate        int    `db:"rate"`
}
