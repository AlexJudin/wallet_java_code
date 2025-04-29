package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/custom_error"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/usecases"
)

var messageError string

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
		log.Errorf("create payment operation: %+v", err)
		messageError = "Переданы некорректные данные о платежной операции."

		returnErr(http.StatusBadRequest, messageError, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &paymentOperation); err != nil {
		log.Errorf("create payment operation: %+v", err)
		messageError = "Не удалось прочитать данные о платежной операции."

		returnErr(http.StatusBadRequest, messageError, w)
		return
	}

	if nameFields, err := checkCreateOperationRequest(paymentOperation); err != nil {
		log.Errorf("create payment operation: %+v", err)
		messageError = fmt.Sprintf("В данных о платежной операции переданы некорректные поля [%s].", nameFields)

		returnErr(http.StatusBadRequest, messageError, w)
		return
	}

	err = h.uc.CreateOperation(&paymentOperation)
	switch {
	case errors.Is(err, custom_error.InsufficientFundsErr):
		log.Errorf("wallet [%s] error: %+v", err)
		messageError = "Недостаточно средств."

		returnErr(http.StatusOK, messageError, w)
		return
	case err != nil:
		log.Errorf("create payment operation: error create payment operation for wallet [%s], operation type [%s], amount [%d]: service is not allowed",
			paymentOperation.WalletId,
			paymentOperation.OperationType,
			paymentOperation.Amount)
		messageError = "Ошибка сервера, не удалось сохранить данные о платежной операции. Попробуйте позже или обратитесь в тех. поддержку."

		returnErr(http.StatusInternalServerError, messageError, w)
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
		messageError = "Не передан идентификатор кошелька, получение баланса невозможно."

		returnErr(http.StatusBadRequest, messageError, w)
		return
	}

	balance, err := h.uc.GetWalletBalanceByUUID(walletUUID)
	if err != nil {
		log.Error("get wallet balance by UUID error: service is not allowed")
		messageError = "Ошибка сервера, не удалось получить баланс. Попробуйте позже или обратитесь в тех. поддержку."

		returnErr(http.StatusInternalServerError, messageError, w)
		return
	}

	respMap := map[string]interface{}{
		"balance": balance,
	}

	resp, err := json.Marshal(respMap)
	if err != nil {
		log.Errorf("get wallet balance by UUID error: %+v", err)
		messageError = "Ошибка сервера. Попробуйте позже или обратитесь в тех. поддержку."

		returnErr(http.StatusInternalServerError, messageError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Errorf("get wallet balance by UUID error: %+v", err)
		messageError = "Сервер недоступен. Попробуйте позже или обратитесь в тех. поддержку."

		returnErr(http.StatusInternalServerError, messageError, w)
	}
}

func returnErr(status int, messageError string, w http.ResponseWriter) {
	message := errResponse{
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
		log.Errorf("get wallet balance by UUID error: %+v", err)
	}
}

func checkCreateOperationRequest(req model.PaymentOperation) (string, error) {
	errorFields := make([]string, 0)

	if req.OperationTypeIsEmpty() {
		errorFields = append(errorFields, "operationType")
	}

	if req.AmountIsNegative() {
		errorFields = append(errorFields, "amount")
	}

	if len(errorFields) == 0 {
		return "", nil
	}

	errorFieldsString := strings.Join(errorFields, ", ")
	errorText := fmt.Sprintf("invalid parameters [%s]", errorFieldsString)

	return errorFieldsString, errors.New(errorText)
}
