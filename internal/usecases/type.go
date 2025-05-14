package usecases

import (
	"context"

	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
	"github.com/AlexJudin/wallet_java_code/internal/model"
)

type Wallet interface {
	CreateOperation(paymentOperation *model.PaymentOperation) error
	GetWalletBalanceByUUID(walletUUID string) (int64, error)

	SetCacheValue(ctx context.Context, walletUUID string, balance int64) error
	GetCacheValue(ctx context.Context, walletUUID string) (int64, error)
}

type Register interface {
	RegisterUser(login string, password string) error
}

type Authorization interface {
	AuthorizationUser(login string, password string) (entity.Tokens, error)
	RefreshToken(refreshToken string) (entity.Tokens, error)
}
