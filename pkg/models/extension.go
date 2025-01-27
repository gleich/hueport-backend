package models

import (
	"time"

	"gorm.io/gorm"
)

type Extension struct {
	gorm.Model
	Name        string
	ExtensionID string
	Themes      int
	LastUpdated time.Time
}
