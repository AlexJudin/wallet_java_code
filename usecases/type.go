package usecases

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(paymentOperation *model.PaymentOperation) error
	GetWalletBalanceByUUID(id string) (int64, error)
}
