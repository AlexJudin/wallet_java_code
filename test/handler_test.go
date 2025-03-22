package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexJudin/wallet_java_code/api"
	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/repository"
	"github.com/AlexJudin/wallet_java_code/usecases"
)

func TestGetWalletBalanceByUUIDWhenOk(t *testing.T) {
	// init config
	cfg, err := config.New()
	if err != nil {
		return
	}

	connStr := cfg.GetDataSourceName()
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		return
	}

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletHandler := api.NewWalletHandler(walletUC)

	req := httptest.NewRequest("GET", "/api/v1/wallets/?WALLET_UUID=ec82ea03-2b53-4258-ba87-a7efae979c43", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(walletHandler.GetWalletBalanceByUUID)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}
}

func TestGetWalletBalanceByUUIDWhenWalletUUIDIsEmpty(t *testing.T) {
	// init config
	cfg, err := config.New()
	if err != nil {
		return
	}

	connStr := cfg.GetDataSourceName()
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		return
	}

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletHandler := api.NewWalletHandler(walletUC)

	req := httptest.NewRequest("GET", "/api/v1/wallets/?WALLET_UUID=", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(walletHandler.GetWalletBalanceByUUID)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code: %d, got %d", http.StatusInternalServerError, status)
	}

	expected := `count missing`
	if responseRecorder.Body.String() != expected {
		t.Errorf("expected body: %s, got %s", expected, responseRecorder.Body.String())
	}
}

func TestGetWalletBalanceByUUIDWhenMissingWalletUUID(t *testing.T) {
	// init config
	cfg, err := config.New()
	if err != nil {
		return
	}

	connStr := cfg.GetDataSourceName()
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		return
	}

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletHandler := api.NewWalletHandler(walletUC)

	req := httptest.NewRequest("GET", "/api/v1/wallets/", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(walletHandler.GetWalletBalanceByUUID)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code: %d, got %d", http.StatusBadRequest, status)
	}

	expected := `count missing`
	if responseRecorder.Body.String() != expected {
		t.Errorf("expected body: %s, got %s", expected, responseRecorder.Body.String())
	}
}
