package repository

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(task *model.PaymentOperation) error
	GetWalletBalanceByUUID(walletUUID string) (int, error)
}
