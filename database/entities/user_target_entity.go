package entities

import "gorm.io/gorm"

type UserTarget struct {
	gorm.Model
	UserID              uint
	UniversityProgramID uint
	Priority            int

	User               User
	UniversityProgram  UniversityProgram
}
