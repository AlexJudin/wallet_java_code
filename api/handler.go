package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/model"
	"github.com/AlexJudin/wallet_java_code/usecases"
)

type WalletHandler struct {
	uc usecases.Wallet
}

func NewWalletHandler(uc usecases.Wallet) WalletHandler {
	return WalletHandler{uc: uc}
}

type errResponse struct {
	Error string `json:"error"`
}

// CreateOperation ... Добавить новую задачу
// @Summary Добавить новую задачу
// @Description Добавить новую задачу
// @Accept json
// @Tags Task
// @Param Body body model.Task true "Параметры задачи"
// @Success 201 {int}    http.StatusCreated
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/task [post]
func (h *WalletHandler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var (
		paymentOperation model.PaymentOperation
		buf              bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Error("create payment operation: error reading body request")

		returnErr(http.StatusBadRequest, err, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &paymentOperation); err != nil {
		log.Error("create payment operation: error unmarshalling body request to payment operation model")

		returnErr(http.StatusBadRequest, err, w)
		return
	}

	err = h.uc.CreateOperation(&paymentOperation)
	if err != nil {
		log.Errorf("http.CreateTask: %+v", err)

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// GetWalletBalanceByUUID ... Получить задачу
// @Summary Получить задачу
// @Description Получить задачу
// @Accept json
// @Tags Task
// @Param id query string true "Идентификатор задачи"
// @Success 200 {object} model.TaskResp
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/task [get]
func (h *WalletHandler) GetWalletBalanceByUUID(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		err := fmt.Errorf("task id is empty")
		log.Errorf("http.GetTask: %+v", err)

		returnErr(http.StatusBadRequest, err, w)
		return
	}

	taskResp, err := h.uc.GetWalletByUUID(taskId)
	if err != nil {
		log.Errorf("http.GetTask: %+v", err)

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	resp, err := json.Marshal(taskResp)
	if err != nil {
		log.Errorf("http.GetTask: %+v", err)

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Errorf("http.GetTask: %+v", err)

		returnErr(http.StatusInternalServerError, err, w)
	}
}

func returnErr(status int, err error, w http.ResponseWriter) {
	message := errResponse{
		Error: err.Error(),
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

	}
}
