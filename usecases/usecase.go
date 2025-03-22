package usecases

import (
	"github.com/AlexJudin/wallet_java_code/model"
	"github.com/AlexJudin/wallet_java_code/repository"
)

var _ Wallet = (*WalletUsecase)(nil)

type WalletUsecase struct {
	DB repository.Wallet
}

func NewWalletUsecase(db repository.Wallet) *WalletUsecase {
	return &WalletUsecase{DB: db}
}

func (t *WalletUsecase) CreateOperation(task *model.PaymentOperation) error {
	err := t.DB.CreateOperation(task)
	if err != nil {
		return err
	}

	return nil
}

func (t *WalletUsecase) GetWalletBalanceByUUID(id string) (*model.PaymentOperation, error) {
	return t.DB.GetWalletBalanceByUUID(id)
}
