package entities

import "gorm.io/gorm"

type University struct {
	gorm.Model

	Name string `json:"name" form:"name" binding:"required"`
	Type string `json:"type" form:"type" binding:"required,oneof=PTN PTS PTK" gorm:"type:enum('PTN','PTS','PTK')"`

	// relations
	Major []UniversityMajor `json:"programs,omitempty"`
}
