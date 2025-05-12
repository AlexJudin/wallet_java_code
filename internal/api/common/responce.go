package common

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/api/entity"
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
