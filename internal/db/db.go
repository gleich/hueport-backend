package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"pkg.mattglei.ch/hueport-scraper/pkg/models"
	"pkg.mattglei.ch/timber"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		timber.Fatal(err, "failed to connect to database")
	}

	err = db.AutoMigrate(&models.Extension{})
	if err != nil {
		timber.Fatal(err, "failed to run migration for extension")
	}

	timber.Done("connected to database")
	return db
}
