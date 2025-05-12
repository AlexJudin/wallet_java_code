package entity

import "github.com/golang-jwt/jwt/v4"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthClaims struct {
	Login string `json:"l"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	Login         string `json:"l"`
	AccessTokenID string `json:"aid"`
	jwt.RegisteredClaims
}
