package entities

import "gorm.io/gorm"

type UserTarget struct {
    gorm.Model

    UserID              uint `json:"userId" form:"userId" binding:"required"`
    UniversityMajorID   uint `json:"universityProgramId" form:"universityProgramId" binding:"required"` // ganti namanya agar jelas foreign key
    Priority            int  `json:"priority" form:"priority" binding:"required,min=1"`

    // relations
    User  User           `json:"user,omitempty"`
    Major UniversityMajor `gorm:"foreignKey:UniversityMajorID;references:ID" json:"universityProgram,omitempty"`
}
