package entities

import (
	"time"

	"gorm.io/gorm"
)

// SubtestResult stores the result of each subtest within a Try Out attempt.
// Used for displaying breakdown per subtest and calculating weighted scores.
type SubtestResult struct {
	gorm.Model

	AttemptID uint `json:"attemptId" gorm:"uniqueIndex:idx_attempt_subtest;not null"`
	SubtestID uint `json:"subtestId" gorm:"uniqueIndex:idx_attempt_subtest;not null"`

	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`

	CorrectCount    int `json:"correctCount" gorm:"default:0"`
	WrongCount      int `json:"wrongCount" gorm:"default:0"`
	UnansweredCount int `json:"unansweredCount" gorm:"default:0"`

	RawScore   *float64 `json:"rawScore" gorm:"type:decimal(10,2)"`   // Weighted score
	FinalScore *float64 `json:"finalScore" gorm:"type:decimal(10,2)"` // Normalized to max_score

	// Relations
	Attempt TryOutAttempt `json:"attempt,omitempty" gorm:"foreignKey:AttemptID"`
	Subtest Subtest       `json:"subtest,omitempty" gorm:"foreignKey:SubtestID"`
}
