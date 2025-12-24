package entities

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	ProgramID       uint
	CreatedByUserID uint

	NameCourse  string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	IsFree      bool `gorm:"default:false"`

	Program  Program
	Subjects []ClassSubject `gorm:"foreignKey:CourseID"`
	Creator  User           `gorm:"foreignKey:CreatedByUserID"`
}
