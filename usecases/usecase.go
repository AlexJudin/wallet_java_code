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

func (t *WalletUsecase) CreateOperation(paymentOperation *model.PaymentOperation) error {
	if paymentOperation.OperationType == model.Withdraw {
		balance, err := t.DB.GetWalletBalanceByUUID(paymentOperation.WalletId)
		if err != nil {
			return err
		}

		if balance < paymentOperation.Amount {
			return model.InsufficientFundsErr
		}

		paymentOperation.Amount = -paymentOperation.Amount
	}

	err := t.DB.CreateOperation(paymentOperation)
	if err != nil {
		return err
	}

	return nil
}

func (t *WalletUsecase) GetWalletBalanceByUUID(walletUUID string) (int64, error) {
	return t.DB.GetWalletBalanceByUUID(walletUUID)
}
