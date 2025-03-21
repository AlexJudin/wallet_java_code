package repository

import (
	"github.com/AlexJudin/wallet_java_code/model"
)

type Wallet interface {
	CreateTask(task *model.Task) (int64, error)
	GetTaskById(id string) (*model.Task, error)
}
