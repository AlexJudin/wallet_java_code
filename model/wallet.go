package model

type PaymentOperationType string

const (
	Deposit  PaymentOperationType = "deposit"
	Withdraw PaymentOperationType = "withdraw"
)

type PaymentOperation struct {
	WalletId      string               `json:"walletId"`
	OperationType PaymentOperationType `json:"operationType"`
	Amount        int64                `json:"amount"`
}
