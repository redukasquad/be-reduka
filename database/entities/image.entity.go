package entities

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	URL    string `gorm:"size:191;not null;uniqueIndex"`
	FileID string `gorm:"size:255;not null"`
}

