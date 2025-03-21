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

func (t *WalletUsecase) CreateOperation(task *model.Wallet) error {
	taskId, err := t.DB.CreateTask(task)
	if err != nil {
		return err
	}

	taskResp := model.NewTaskResp(taskId)

	return nil
}

func (t *WalletUsecase) GetWalletByUUID(id string) (int, error) {
	return t.DB.GetTaskById(id)
}
