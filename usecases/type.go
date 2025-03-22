package usecases

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(task *model.PaymentOperation) error
	GetWalletBalanceByUUID(id string) (int, error)
}
