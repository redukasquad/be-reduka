package entities

import "gorm.io/gorm"

type Program struct {
	gorm.Model
	ProgramName string `gorm:"type:varchar(100);uniqueIndex"`
	Description string `gorm:"type:text"`
	ImageProgram string `gorm:"type:text"`

	Courses []Course
}
