package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

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

		custom_error.ReturnHTTPErr(http.StatusBadRequest, messageError, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		log.Errorf("authorization user error: %+v", err)
		messageError = "Не удалось прочитать логин/пароль."

		custom_error.ReturnHTTPErr(http.StatusBadRequest, messageError, w)
		return
	}

	tokens, err := h.uc.AuthorizationUser(user.Login, user.Password)
	switch {
	case errors.Is(err, custom_error.ErrNotFound):
		log.Errorf("authorization user error: %+v", err)
		messageError = "Пользователь не найден."

		custom_error.ReturnHTTPErr(http.StatusNotFound, messageError, w)
		return
	case errors.Is(err, custom_error.ErrIncorrectPassword):
		log.Errorf("authorization user error: %+v", err)
		messageError = "Некорректный пароль."

		custom_error.ReturnHTTPErr(http.StatusForbidden, messageError, w)
		return
	case err != nil:
		log.Errorf("authorization user error: %+v", err)
		messageError = "Ошибка сервера, не удалось авторизовать пользователя. Попробуйте позже или обратитесь в тех. поддержку."

		custom_error.ReturnHTTPErr(http.StatusInternalServerError, messageError, w)
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
