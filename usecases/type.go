package usecases

import (
	"github.com/AlexJudin/go_final_project/model"
)

type Wallet interface {
	CreateTask(task *model.Task, today bool) (*model.TaskResp, error)
	GetTaskById(id string) (*model.Task, error)
}
