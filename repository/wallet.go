package repository

import (
	"database/sql"

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

func (r *WalletRepo) CreateOperation(task *model.PaymentOperation) error {
	_, err := r.Db.Exec(SQLCreateTask)
	if err != nil {
		log.Debugf("Database.CreateTask: %+v", err)

		return err
	}

	return nil
}

func (r *WalletRepo) GetWalletByUUID(id string) (*model.PaymentOperation, error) {
	var task model.PaymentOperation

	res, err := r.Db.Query(SQLGetTaskById, id)
	if err != nil {
		log.Debugf("Database.GetTaskById: %+v", err)

		return nil, err
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Debugf("Database.GetTaskById: %+v", err)

			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	/*
		if task.Id == "" {
			err = fmt.Errorf("task id %s not found", id)
			log.Debugf("Database.GetTaskById: %+v", err)

			return nil, err
		}
	*/

	return &task, nil
}
