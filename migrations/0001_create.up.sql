CREATE TABLE IF NOT EXISTS wallets
(
    id             BIGINT NOT NULL PRIMARY KEY,
    wallet_guid    uuid NOT NULL,
    operation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    operation_type TEXT NOT NULL,
    amount         INT NOT NULL
);

COMMENT ON TABLE wallets IS 'таблица платежных операций по кошелькам';
COMMENT ON COLUMN wallets.id IS 'уникальный идентификатор';
COMMENT ON COLUMN wallets.wallet_guid IS 'ГУИД кошелька';
COMMENT ON COLUMN wallets.operation_date IS 'дата проведения платежной операции';
COMMENT ON COLUMN wallets.operation_type IS 'тип платежной операции';
COMMENT ON COLUMN wallets.amount IS 'сумма по платежной операции';
