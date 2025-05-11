package custom_error

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	ErrNotFound          = errors.New("entity not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type errHTTPResponse struct {
	Error string `json:"error"`
}

func ReturnHTTPErr(status int, messageError string, w http.ResponseWriter) {
	message := errHTTPResponse{
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
