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
	WalletId      string               `json:"walletId" gorm:"not null"`
	OperationType PaymentOperationType `json:"operationType" gorm:"not null"`
	Amount        int64                `json:"amount" gorm:"default 0"`
}
