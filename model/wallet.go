package model

type OperationType string

const (
	Deposit  OperationType = "deposit"
	Withdraw OperationType = "withdraw"
)

type Operation struct {
	WalletId      string        `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        int64         `json:"amount"`
}
