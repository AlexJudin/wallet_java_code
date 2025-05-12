package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/api/common"
	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/usecases"
)

var messageError string

type AuthHandler struct {
	uc usecases.Authorization
}

func NewAuthHandler(uc usecases.Authorization) AuthHandler {
	return AuthHandler{uc: uc}
}

func (h *AuthHandler) AuthorizationUser(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Errorf("authorization user error: %+v", err)
		messageError = "Переданы некорректные логин/пароль."

		common.ApiError(http.StatusBadRequest, messageError, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		log.Errorf("authorization user error: %+v", err)
		messageError = "Не удалось прочитать логин/пароль."

		common.ApiError(http.StatusBadRequest, messageError, w)
		return
	}

	tokens, err := h.uc.AuthorizationUser(user.Login, user.Password)
	switch {
	case errors.Is(err, custom_error.ErrNotFound):
		log.Errorf("authorization user error: %+v", err)
		messageError = "Пользователь не найден."

		common.ApiError(http.StatusNotFound, messageError, w)
		return
	case errors.Is(err, custom_error.ErrIncorrectPassword):
		log.Errorf("authorization user error: %+v", err)
		messageError = "Некорректный пароль."

		common.ApiError(http.StatusForbidden, messageError, w)
		return
	case err != nil:
		log.Errorf("authorization user error: %+v", err)
		messageError = "Ошибка сервера, не удалось авторизовать пользователя. Попробуйте позже или обратитесь в тех. поддержку."

		common.ApiError(http.StatusInternalServerError, messageError, w)
		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		HttpOnly: true,
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		log.Error("refresh token not found")
		messageError = "refresh token не найден"

		common.ApiError(http.StatusUnauthorized, messageError, w)
		return
	}

	tokens, err := h.uc.RefreshToken(refreshToken.Value)
	if err != nil {
		if errors.Is(err, custom_error.ErrNotFound) {
			log.Error("token not found in storage")
			messageError = "Токен не найден"

			common.ApiError(http.StatusNotFound, messageError, w)
			return
		}

		log.Error("cannot refresh token")
		messageError = "Не удалось обновить токен"

		common.ApiError(http.StatusInternalServerError, messageError, w)
		return
	}

	accessTokenCookie := http.Cookie{
		Name:     "accessToken",
		Value:    tokens.AccessToken,
		HttpOnly: true,
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	}
	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	w.WriteHeader(http.StatusOK)
}
