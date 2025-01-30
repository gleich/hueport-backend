package models

import (
	"time"

	"gorm.io/gorm"
)

type Extension struct {
	gorm.Model
	Name        string
	ExtensionID string `gorm:"uniqueIndex"`
	LastUpdated time.Time
	Themes      []Theme `gorm:"foreignKey:ExtensionID;references:ExtensionID"`
}
