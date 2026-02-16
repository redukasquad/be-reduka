package entities

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model

	ProgramID       uint `json:"programId" form:"programId" binding:"required"`
	CreatedByUserID uint `json:"createdByUserId" form:"createdByUserId" binding:"required"`

	NameCourse        string    `json:"nameCourse" form:"nameCourse" binding:"required"`
	Description       string    `json:"description" form:"description"`
	StartDate         time.Time `json:"startDate" form:"startDate" binding:"required"`
	EndDate           time.Time `json:"endDate" form:"endDate" binding:"required"`
	IsFree            bool      `json:"isFree" form:"isFree" gorm:"default:false"`
	WhatsappGroupLink string    `json:"whatsappGroupLink" form:"whatsappGroupLink"`
	image 	 		  string    `json:"image" form:"image"`

	// relations
	Program   Program                `json:"program,omitempty"`
	Subjects  []ClassSubject         `json:"subjects,omitempty" gorm:"foreignKey:CourseID"`
	Creator   User                   `json:"creator,omitempty" gorm:"foreignKey:CreatedByUserID"`
	Questions []RegistrationQuestion `gorm:"foreignKey:CourseID"`
}
