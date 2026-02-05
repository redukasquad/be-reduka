package entities

import "gorm.io/gorm"

// DifficultyLevel represents the difficulty classification for IRT scoring.
type DifficultyLevel string

const (
	DifficultyEasy   DifficultyLevel = "easy"
	DifficultyMedium DifficultyLevel = "medium"
	DifficultyHard   DifficultyLevel = "hard"
)

// TryOutQuestion represents a question in a Try Out subtest.
// Options A-E are stored as columns (denormalized) since they are always fixed 5 options.
type TryOutQuestion struct {
	gorm.Model

	TryOutPackageID uint `json:"tryOutPackageId" gorm:"index;not null"`
	SubtestID       uint `json:"subtestId" gorm:"index;not null"`

	QuestionText string `json:"questionText" gorm:"type:text;not null"`
	ImageURL     string `json:"imageUrl" gorm:"size:500"`     // Optional, ImageKit URL
	Explanation  string `json:"explanation" gorm:"type:text"` // Pembahasan

	DifficultyLevel DifficultyLevel `json:"difficultyLevel" gorm:"size:10;not null"` // easy, medium, hard
	OrderNumber     int             `json:"orderNumber" gorm:"not null"`

	// Options (denormalized - always A-E)
	OptionA string `json:"optionA" gorm:"type:text;not null"`
	OptionB string `json:"optionB" gorm:"type:text;not null"`
	OptionC string `json:"optionC" gorm:"type:text;not null"`
	OptionD string `json:"optionD" gorm:"type:text;not null"`
	OptionE string `json:"optionE" gorm:"type:text;not null"`

	CorrectOption   string `json:"correctOption" gorm:"size:1;not null"` // A, B, C, D, or E
	CreatedByUserID uint   `json:"createdByUserId"`

	// Relations
	TryOutPackage TryOut `json:"tryOutPackage,omitempty" gorm:"foreignKey:TryOutPackageID"`
	Subtest       Subtest       `json:"subtest,omitempty" gorm:"foreignKey:SubtestID"`
	Creator       User          `json:"creator,omitempty" gorm:"foreignKey:CreatedByUserID"`
}
