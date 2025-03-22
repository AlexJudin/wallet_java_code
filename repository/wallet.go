package repository

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/model"
)

var _ Wallet = (*WalletRepo)(nil)

type WalletRepo struct {
	Db *sql.DB
}

func NewWalletRepo(db *sql.DB) *WalletRepo {
	return &WalletRepo{Db: db}
}

func (r *WalletRepo) CreateOperation(paymentOperation *model.PaymentOperation) error {
	log.Infof("start saving payment operation for wallet [%s]", paymentOperation.WalletId)

	sqlText, args, err := sq.Insert("wallets").
		Columns("wallet_guid", "operation_type", "amount").
		Values(paymentOperation.WalletId, paymentOperation.OperationType, paymentOperation.Amount).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Debugf("unable to build INSERT query: %+v", err)
		return err
	}

	log.Infof("executing SQL: %s", sqlText)

	_, err = r.Db.Exec(sqlText, args...)
	if err != nil {
		log.Debugf("error create payment operation: %+v", err)
		return err
	}

	return nil
}

func (r *WalletRepo) GetWalletBalanceByUUID(walletUUID string) (int, error) {
	log.Infof("start getting balance for wallet [%s]", walletUUID)

	var balance int

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

	return balance, nil
}
