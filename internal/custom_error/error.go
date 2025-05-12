package custom_error

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
)

var (
	ErrNotFound          = errors.New("entity not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func ApiError(status int, messageError string, w http.ResponseWriter) {
	message := entity.ApiError{
		Error: messageError,
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		status = http.StatusInternalServerError
		messageJson = []byte("{\"error\":\"" + err.Error() + "\"}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(messageJson)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Errorf("unknow server error: %+v", err)
	}
}
