package entities

import "gorm.io/gorm"

type UniversityProgram struct {
	gorm.Model
	UniversityID uint
	Name         string
	PassingGrade float64

	University University
	Targets    []UserTarget
}
