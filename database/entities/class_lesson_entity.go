package entities

import (
	"time"

	"gorm.io/gorm"
)

type ClassLesson struct {
	gorm.Model
	SubjectID       uint
	CreatedByUserID uint

	Title       string
	Description string
	LessonOrder int

	StartTime *time.Time
	EndTime   *time.Time

	Subject ClassSubject `gorm:"foreignKey:SubjectID"`
	Creator User `gorm:"foreignKey:CreatedByUserID"`
	Resources []ClassLessonResource `gorm:"foreignKey:ClassLessonID"`
}