package repository

import (
	"gorm.io/gorm"

	"github.com/AlexJudin/wallet_java_code/internal/model"
)

var _ Register = (*RegisterRepo)(nil)

type RegisterRepo struct {
	Db *gorm.DB
}

func NewRegisterRepo(db *gorm.DB) *RegisterRepo {
	return &RegisterRepo{Db: db}
}

func (r *RegisterRepo) SaveUser(login string, password string) error {
	return nil
}

func (r *RegisterRepo) GetUserByLogin(login string) (model.User, error) {

}
