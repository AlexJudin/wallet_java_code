package usecases

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(task *model.Wallet) error
	GetWalletByUUID(id string) (int, error)
}
