package entities

import "gorm.io/gorm"

type University struct {
	gorm.Model
	Name string
	Type string `gorm:"type:enum('PTN','PTS','PTK')"`

	Programs []UniversityProgram
}