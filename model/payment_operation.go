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
	ID            uint                 `gorm:"primarykey" json:"-"`
	CreatedAt     time.Time            `json:"-"`
	WalletId      string               `json:"walletId"`
	OperationType PaymentOperationType `json:"operationType"`
	Amount        int64                `json:"amount"`
}

func (p PaymentOperation) OperationTypeIsEmpty() bool {
	return p.OperationType == ""
}

func (p PaymentOperation) AmountIs() bool {
	return p.Amount < 0
}
