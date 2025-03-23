package test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/AlexJudin/wallet_java_code/api"
	"github.com/AlexJudin/wallet_java_code/repository"
	"github.com/AlexJudin/wallet_java_code/usecases"
)

var (
	walletTest WalletTest
)

type WalletTest struct {
	db            *gorm.DB
	walletHandler api.WalletHandler
}

func TestMain(m *testing.M) {
	if err := initialize(); err != nil {
		log.Panicf("error initializing test environment: %+v", err)
	}
	os.Exit(m.Run())
}

func initialize() error {
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		return err
	}
	walletTest.db = db

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletTest.walletHandler = api.NewWalletHandler(walletUC)

	return nil
}
