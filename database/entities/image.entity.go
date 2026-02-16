package entities

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	URL   string `gorm:"type:text;not null;uniqueIndex"`
	Fileid string `gorm:"type:varchar(255);not null"`
}
