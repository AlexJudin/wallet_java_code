package usecases

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateOperation(task *model.Operation) error
	GetWalletByUUID(id string) (*model.Operation, error)
}
