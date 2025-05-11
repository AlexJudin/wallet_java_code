package usecases

import (
	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
	"github.com/AlexJudin/wallet_java_code/internal/service"
)

var _ Register = (*RegisterUsecase)(nil)

type RegisterUsecase struct {
	DB          repository.User
	ServiceAuth service.AuthService
}

func NewRegisterUsecase(db repository.User, serviceAuth service.AuthService) *RegisterUsecase {
	return &RegisterUsecase{
		DB:          db,
		ServiceAuth: serviceAuth,
	}
}

func (u *RegisterUsecase) RegisterUser(login string, password string) error {
	user, err := u.DB.GetUserByLogin(login)
	if err != nil {
		return err
	}

	if user.IsAlreadyExist() {
		return custom_error.ErrUserAlreadyExists
	}

	newUser := model.User{
		Login: login,
		Hash:  u.ServiceAuth.GenerateHashPassword(password),
	}

	err = u.DB.SaveUser(newUser)
	if err != nil {
		return err
	}

	return nil
}
