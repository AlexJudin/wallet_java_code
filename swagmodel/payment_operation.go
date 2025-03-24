package swagmodel

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type PaymentOperation struct {
	WalletId      string                     `json:"walletId"`
	OperationType model.PaymentOperationType `json:"operationType"`
	Amount        int64                      `json:"amount"`
}
