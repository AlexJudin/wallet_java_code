package repository

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(paymentOperation *model.PaymentOperation) error
	GetWalletBalanceByUUID(walletUUID string) (int, error)
}
