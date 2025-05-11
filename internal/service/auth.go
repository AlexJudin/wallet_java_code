package service

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/internal/model"
)

type AuthService struct {
	Config *config.Сonfig
}

func NewAuthService(cfg *config.Сonfig) AuthService {
	return AuthService{Config: cfg}
}

func (s AuthService) GenerateHashPassword(password string) string {
	passwordBytes := []byte(password)
	sha512Hasher := sha512.New()

	passwordBytes = append(passwordBytes, s.Config.PasswordSalt...)
	sha512Hasher.Write(passwordBytes)

	hashedPasswordBytes := sha512Hasher.Sum(nil)
	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

func (s AuthService) GenerateTokens(login string) (model.Tokens, error) {
	accessTokenID := uuid.NewString()
	accessToken, err := s.generateAccessToken(login)
	if err != nil {
		return model.Tokens{}, err
	}

	refreshToken, err := s.generateRefreshToken(login, accessTokenID)
	if err != nil {
		return model.Tokens{}, err
	}

	return model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s AuthService) generateAccessToken(login string) (string, error) {

}

func (s AuthService) generateRefreshToken(login string, accessTokenID string) (string, error) {

}
