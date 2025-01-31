package models

import (
	"time"

	"gorm.io/gorm"
)

type Extension struct {
	ExtensionID string `gorm:"primaryKey"`
	Name        string
	LastUpdated time.Time
	Themes      []Theme
}

type Theme struct {
	gorm.Model
	Name          string
	ExtensionID   string
	Foreground    string
	Background    string
	BrightWhite   string
	White         string
	BrightBlack   string
	Black         string
	BrightBlue    string
	Blue          string
	BrightGreen   string
	Green         string
	BrightCyan    string
	Cyan          string
	BrightRed     string
	Red           string
	BrightMagenta string
	Magenta       string
	BrightYellow  string
	Yellow        string
}
