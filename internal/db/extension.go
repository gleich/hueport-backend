package db

import (
	"time"

	"gorm.io/gorm"
)

type Extension struct {
	gorm.Model
	Name        string
	ExtensionID string
	LastUpdated time.Time
}
