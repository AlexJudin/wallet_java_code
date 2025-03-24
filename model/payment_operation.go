package model

import (
	"time"
)

type PaymentOperationType string

const (
	Deposit  PaymentOperationType = "deposit"
	Withdraw PaymentOperationType = "withdraw"
)

type PaymentOperation struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	WalletId      string               `json:"walletId"`
	OperationType PaymentOperationType `json:"operationType"`
	Amount        int64                `json:"amount"`
}
