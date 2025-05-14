package usecases

import (
	"context"
	"strconv"

	"github.com/AlexJudin/wallet_java_code/internal/cache"
	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
)

var _ Wallet = (*WalletUsecase)(nil)

type WalletUsecase struct {
	DB    repository.Wallet
	Cache cache.Client
}

func NewWalletUsecase(db repository.Wallet, cache cache.Client) *WalletUsecase {
	return &WalletUsecase{
		DB:    db,
		Cache: cache,
	}
}

func (t *WalletUsecase) CreateOperation(paymentOperation *model.PaymentOperation) error {
	if paymentOperation.OperationType == model.Withdraw {
		balance, err := t.DB.GetWalletBalanceByUUID(paymentOperation.WalletId)
		if err != nil {
			return err
		}

		if balance < paymentOperation.Amount {
			return custom_error.ErrInsufficientFunds
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

func (t *WalletUsecase) SetCacheValue(ctx context.Context, walletUUID string, balance int64) error {
	return t.Cache.SetValue(ctx, walletUUID, balance)
}

func (t *WalletUsecase) GetCacheValue(ctx context.Context, walletUUID string) (int64, error) {
	result, err := t.Cache.GetValue(ctx, walletUUID)
	if err != nil {
		return 0, err
	}

	balance, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}

	return int64(balance), nil
}
