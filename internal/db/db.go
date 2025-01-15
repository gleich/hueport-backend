package db

import (
	"os"

	"github.com/gleich/lumber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		lumber.Fatal(err, "failed to connect to database")
	}

	err = db.AutoMigrate(&Extension{})
	if err != nil {
		lumber.Fatal(err, "failed to run migration for extension")
	}

	lumber.Done("connected to database")
	return db
}
