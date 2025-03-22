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

	paymentOperation = new(model.PaymentOperation)

	err := r.Db.Create(paymentOperation).Error
	if err != nil {
		log.Debugf("error create payment operation: %+v", err)
		return err
	}

	return nil
}

func (r *WalletRepo) GetWalletBalanceByUUID(walletUUID string) (int, error) {
	log.Infof("start getting balance for wallet [%s]", walletUUID)

	var balance int

	/*
		sqlText, args, err := sq.Select("SUM(amount) AS amount").
			From("wallets").
			Where(sq.Eq{"wallet_guid": walletUUID}).
			PlaceholderFormat(sq.Dollar).
			GroupBy("wallet_guid").
			ToSql()
		if err != nil {
			log.Debugf("unable to build SELECT query: %+v", err)
			return 0, err
		}

		log.Infof("executing SQL: %s", sqlText)

		res, err := r.Db.Query(sqlText, args...)
		if err != nil {
			log.Debugf("error get balance for wallet [%s]: %+v", walletUUID, err)
			return 0, err
		}
		defer res.Close()

		if res.Next() {
			err = res.Scan(&balance)
			if err != nil {
				log.Debug(err)
				return 0, err
			}
		}

		if err = res.Err(); err != nil {
			return 0, err
		}
	*/

	return balance, nil
}
