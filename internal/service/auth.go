package service

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
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

func (s AuthService) GenerateTokens(login string) (entity.Tokens, error) {
	accessTokenID := uuid.NewString()
	accessToken, err := s.generateAccessToken(login)
	if err != nil {
		return entity.Tokens{}, err
	}

	refreshToken, err := s.generateRefreshToken(login, accessTokenID)
	if err != nil {
		return entity.Tokens{}, err
	}

	// Добавить сохранение токена в БД

	return entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s AuthService) VerifyUser(token string) (string, error) {
	claims := &entity.AuthClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return s.Config.TokenSalt, nil
	})

	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("incorrect token: %+v", err)
	}

	return claims.Login, nil
}

func (s AuthService) generateAccessToken(login string) (string, error) {
	now := time.Now()
	claims := entity.AuthClaims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.Config.TokenSalt) //переделать на []byte
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s AuthService) generateRefreshToken(login string, accessTokenID string) (string, error) {
	now := time.Now()
	claims := entity.RefreshTokenClaims{
		Login:         login,
		AccessTokenID: accessTokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.Config.TokenSalt)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
