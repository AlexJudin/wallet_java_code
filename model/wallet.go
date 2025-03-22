package model

import "gorm.io/gorm"

type PaymentOperationType string

const (
	Deposit  PaymentOperationType = "deposit"
	Withdraw PaymentOperationType = "withdraw"
)

type PaymentOperation struct {
	gorm.Model
	ID            uint                 `json:"id" gorm:"primaryKey"`
	WalletId      string               `json:"walletId"`
	OperationType PaymentOperationType `json:"operationType"`
	Amount        int64                `json:"amount"`
}
