package entities

import "gorm.io/gorm"

type UserTarget struct {
    gorm.Model

    UserID              uint `json:"userId" form:"userId" binding:"required"`
    UniversityMajorID   uint `json:"universityProgramId" form:"universityProgramId" binding:"required"` // ganti namanya agar jelas foreign key
    Priority            int  `json:"priority" form:"priority" binding:"required,min=1"`

<<<<<<< HEAD
    // relations
    User  User           `json:"user,omitempty"`
    Major UniversityMajor `gorm:"foreignKey:UniversityMajorID;references:ID" json:"universityProgram,omitempty"`
=======
	// relations
	User              User              `json:"user,omitempty"`
	Major 						UniversityMajor 	`json:"universityProgram,omitempty" gorm:"foreignKey:UniversityProgramID"`
>>>>>>> 48a3d2b (fix: GORM foreign key tags and route param conflict)
}
