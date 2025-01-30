package models

import "gorm.io/gorm"

type Theme struct {
	gorm.Model
	Name          string
	ExtensionID   string
	Extension     Extension `gorm:"foreignKey:ExtensionID;references:ExtensionID"`
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
