package repository

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(task *model.Wallet) (int64, error)
	GetWalletByUUID(id string) (*model.Wallet, error)
}
