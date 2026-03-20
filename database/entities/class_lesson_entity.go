package entities

import (
	"time"

	"gorm.io/gorm"
)

type Lesson struct {
	gorm.Model

	ClassID         uint `json:"classId" form:"classId" binding:"required"`
	CreatedByUserID uint `json:"createdByUserId" form:"createdByUserId" binding:"required"`

	Title       string `json:"title" form:"title" binding:"required"`
	Description string `json:"description" form:"description"`
	LessonOrder int    `json:"lessonOrder" form:"lessonOrder" binding:"required,min=1"`

	StartTime *time.Time `json:"startTime,omitempty" form:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty" form:"endTime"`

	Class     Class            `json:"class,omitempty" gorm:"foreignKey:ClassID"`
	Creator   User             `json:"creator,omitempty" gorm:"foreignKey:CreatedByUserID"`
	Resources []LessonResource `json:"resources,omitempty" gorm:"foreignKey:LessonID"`
}
