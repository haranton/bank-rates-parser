CREATE TABLE IF NOT EXISTS bank_rates (
    id BIGSERIAL PRIMARY KEY,
    bank_name TEXT NOT NULL UNIQUE,
    deposit_name TEXT NOT NULL,
    rate INT NOT NULL
);