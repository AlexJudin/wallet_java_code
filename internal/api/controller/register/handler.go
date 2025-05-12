package register

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

type RegisterHandler struct {
	uc usecases.Register
}

func NewRegisterHandler(uc usecases.Register) RegisterHandler {
	return RegisterHandler{uc: uc}
}

func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Errorf("register user error: %+v", err)
		messageError = "Переданы некорректные логин/пароль."

		common.ApiError(http.StatusBadRequest, messageError, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		log.Errorf("register user error: %+v", err)
		messageError = "Не удалось прочитать логин/пароль."

		common.ApiError(http.StatusBadRequest, messageError, w)
		return
	}

	err = h.uc.RegisterUser(user.Login, user.Password)
	switch {
	case errors.Is(err, custom_error.ErrUserAlreadyExists):
		log.Errorf("register user error: %+v", err)
		messageError = "Пользователь уже зарегистрирован."

		common.ApiError(http.StatusConflict, messageError, w)
		return
	case err != nil:
		log.Errorf("register user error: %+v", err)
		messageError = "Ошибка сервера, не удалось зарегистрировать пользователя. Попробуйте позже или обратитесь в тех. поддержку."

		common.ApiError(http.StatusInternalServerError, messageError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
