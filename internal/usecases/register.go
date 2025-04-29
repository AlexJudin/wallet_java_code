package usecases

import "github.com/AlexJudin/wallet_java_code/internal/repository"

var _ Register = (*RegisterUsecase)(nil)

type RegisterUsecase struct {
	DB repository.Register
}

func NewRegisterUsecase(db repository.Register) *RegisterUsecase {
	return &RegisterUsecase{DB: db}
}

func (u *RegisterUsecase) RegisterUser(login string, password string) error {
	return nil
}
