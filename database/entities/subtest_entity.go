package entities

import "gorm.io/gorm"

// Subtest represents the 7 fixed subtests for Try Out (PU, PBM, PPU, PK, LBI, LBE, PM).
// This is master data that should be seeded once.
type Subtest struct {
	gorm.Model

	Code             string  `json:"code" gorm:"uniqueIndex;size:10;not null"` // PU, PBM, PPU, etc.
	Name             string  `json:"name" gorm:"size:100;not null"`
	QuestionCount    int     `json:"questionCount" gorm:"not null"`
	TimeLimitSeconds int     `json:"timeLimitSeconds" gorm:"not null"`
	MaxScore         float64 `json:"maxScore" gorm:"type:decimal(10,2);not null"`
}
