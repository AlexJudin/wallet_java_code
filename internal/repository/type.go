package repository

import (
	"github.com/AlexJudin/wallet_java_code/internal/model"
)

type Wallet interface {
	CreateOperation(paymentOperation *model.PaymentOperation) error
	GetWalletBalanceByUUID(walletUUID string) (int64, error)
}

type Register interface {
	GetUserByLogin(login string) (model.User, error)
	SaveUser(user model.User) error
}
