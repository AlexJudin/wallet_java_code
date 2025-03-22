package repository

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/AlexJudin/wallet_java_code/model"
)

var _ Wallet = (*WalletRepo)(nil)

type WalletRepo struct {
	Db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{Db: db}
}

func (r *WalletRepo) CreateOperation(paymentOperation *model.PaymentOperation) error {
	log.Infof("start saving payment operation for wallet [%s]", paymentOperation.WalletId)

	err := r.Db.Create(&paymentOperation).Error
	if err != nil {
		log.Debugf("error create payment operation: %+v", err)
		return err
	}

	return nil
}

func (r *WalletRepo) GetWalletBalanceByUUID(walletUUID string) (int, error) {
	log.Infof("start getting balance for wallet [%s]", walletUUID)

	result := struct {
		Balance int
	}{}

	err := r.Db.Model(&model.PaymentOperation{}).
		Select("SUM(amount) AS balance").
		Where("wallet_id = ?", walletUUID).
		Group("wallet_id").
		Find(&result).Error
	if err != nil {
		log.Debugf("error getting balance for wallet [%s]: %+v", walletUUID, err)
		return 0, err
	}

	return result.Balance, nil
}
