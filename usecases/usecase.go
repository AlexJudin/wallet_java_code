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

func (t *WalletUsecase) CreateTask(task *model.Task, pastDay bool) (*model.TaskResp, error) {
	taskId, err := t.DB.CreateTask(task)
	if err != nil {
		return nil, err
	}

	taskResp := model.NewTaskResp(taskId)

	return taskResp, nil
}

func (t *WalletUsecase) GetTaskById(id string) (*model.Task, error) {
	return t.DB.GetTaskById(id)
}
