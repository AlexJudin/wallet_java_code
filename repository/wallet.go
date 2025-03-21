package repository

import (
	"fmt"

	//sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/model"
)

var _ Wallet = (*WalletRepo)(nil)

type WalletRepo struct {
	Db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
	return &WalletRepo{Db: db}
}

func (r *WalletRepo) CreateTask(task *model.Task) (int64, error) {
	res, err := r.Db.Exec(SQLCreateTask, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Debugf("Database.CreateTask: %+v", err)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Debugf("Database.CreateTask: %+v", err)

		return 0, err
	}

	return id, nil
}

func (r *WalletRepo) GetTaskById(id string) (*model.Task, error) {
	var task model.Task

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

	if task.Id == "" {
		err = fmt.Errorf("task id %s not found", id)
		log.Debugf("Database.GetTaskById: %+v", err)

		return nil, err
	}

	return &task, nil
}
