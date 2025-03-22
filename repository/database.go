package repository

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/AlexJudin/wallet_java_code/model"
)

func ConnectDB(connStr string) (*gorm.DB, error) {
	log.Info("Start connection to database")

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Info("Connected to database")

	log.Info("Running migration")
	db.AutoMigrate(&model.PaymentOperation{})

	return db, nil
}
