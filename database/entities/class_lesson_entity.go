package entities

import (
	"time"

	"gorm.io/gorm"
)

type ClassLesson struct {
	gorm.Model

	SubjectID       uint `json:"subjectId" form:"subjectId" binding:"required"`
	CreatedByUserID uint `json:"createdByUserId" form:"createdByUserId" binding:"required"`

	Title       string `json:"title" form:"title" binding:"required"`
	Description string `json:"description" form:"description"`
	LessonOrder int    `json:"lessonOrder" form:"lessonOrder" binding:"required,min=1"`

	StartTime *time.Time `json:"startTime,omitempty" form:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" form:"endTime"`

	Subject   ClassSubject          `json:"subject,omitempty" gorm:"foreignKey:SubjectID"`
	Creator   User                  `json:"creator,omitempty" gorm:"foreignKey:CreatedByUserID"`
	Resources []ClassLessonResource `json:"resources,omitempty" gorm:"foreignKey:ClassLessonID"`
}