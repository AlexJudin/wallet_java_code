package usecases

import (
	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
	"github.com/AlexJudin/wallet_java_code/internal/service"
)

var _ Authorization = (*AuthUsecase)(nil)

type AuthUsecase struct {
	DB          repository.User
	ServiceAuth service.AuthService
}

func NewAuthUsecase(db repository.User, serviceAuth service.AuthService) *AuthUsecase {
	return &AuthUsecase{
		DB:          db,
		ServiceAuth: serviceAuth,
	}
}

func (u *AuthUsecase) AuthorizationUser(login string, password string) (entity.Tokens, error) {
	user, err := u.DB.GetUserByLogin(login)
	if err != nil {
		return entity.Tokens{}, err
	}

	if user.IsNotFound() {
		return entity.Tokens{}, custom_error.ErrNotFound
	}

	passwordHash := u.ServiceAuth.GenerateHashPassword(password)
	if user.Hash != passwordHash {
		return entity.Tokens{}, custom_error.ErrIncorrectPassword
	}

	return u.ServiceAuth.GenerateTokens(login)
}
