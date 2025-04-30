package usecases

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
)

var _ Register = (*RegisterUsecase)(nil)

type RegisterUsecase struct {
	DB repository.Register
}

func NewRegisterUsecase(db repository.Register) *RegisterUsecase {
	return &RegisterUsecase{DB: db}
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
		Hash:  generateHashPassword(password),
	}

	err = u.DB.SaveUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

func generateHashPassword(password string) string {
	passwordBytes := []byte(password)
	sha512Hasher := sha512.New()

	sha512Hasher.Write(passwordBytes)

	hashedPasswordBytes := sha512Hasher.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}
