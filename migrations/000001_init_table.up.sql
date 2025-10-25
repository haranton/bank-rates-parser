CREATE TABLE IF NOT EXISTS bank_rates (
    id BIGSERIAL PRIMARY KEY,
    bank_name TEXT NOT NULL UNIQUE,
    deposit_name TEXT NOT NULL,
    rate NUMERIC(10,2) NOT NULL  -- например две цифры после запятой
);