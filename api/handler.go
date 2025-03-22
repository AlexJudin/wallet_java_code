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

// CreateOperation ... Добавить новую платежную операцию
// @Summary Добавить новую платежную операцию
// @Description Добавить новую платежную операцию по кошельку
// @Accept json
// @Tags wallet
// @Param Body body model.PaymentOperation true "Параметры операции"
// @Success 201 {int}    http.StatusCreated
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/v1/wallet [post]
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
		log.Errorf("create payment operation: error create payment operation for wallet [%s], operation type [%s], amount [%d]: service is not allowed",
			paymentOperation.WalletId,
			paymentOperation.OperationType,
			paymentOperation.Amount)

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// GetWalletBalanceByUUID ... Получить баланс по кошельку
// @Summary Получить баланс по кошельку
// @Description Получить баланс по кошельку
// @Accept json
// @Tags wallet
// @Param WALLET_UUID query string true "Идентификатор кошелька"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/v1/wallets/ [get]
func (h *WalletHandler) GetWalletBalanceByUUID(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.FormValue("WALLET_UUID")
	if walletUUID == "" {
		err := fmt.Errorf("wallet UUID is empty")
		log.Errorf("get wallet balance by UUID error: %+v", err)

		returnErr(http.StatusBadRequest, err, w)
		return
	}

	balance, err := h.uc.GetWalletBalanceByUUID(walletUUID)
	if err != nil {
		log.Error("get wallet balance by UUID error: service is not allowed")

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	respMap := map[string]interface{}{
		"balance": balance,
	}

	resp, err := json.Marshal(respMap)
	if err != nil {
		log.Errorf("http.GetTask: %+v", err)

		returnErr(http.StatusInternalServerError, err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Errorf("get wallet balance by UUID error: %+v", err)

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
		log.Errorf("get wallet balance by UUID error: %+v", err)
	}
}
