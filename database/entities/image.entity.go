package entities

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	URL    string `gorm:"type:varchar(2048);not null"`
	Fileid string `gorm:"type:varchar(255);not null"`
}
