package entities

import (
	"time"

	"gorm.io/gorm"
)

// UserTryOutAnswer stores each answer submitted by a user during a Try Out attempt.
// Used for displaying which questions were correct/wrong and their explanations.
type UserTryOutAnswer struct {
	gorm.Model

	AttemptID  uint `json:"attemptId" gorm:"uniqueIndex:idx_attempt_question;not null"`
	QuestionID uint `json:"questionId" gorm:"uniqueIndex:idx_attempt_question;not null"`

	SelectedOption *string    `json:"selectedOption" gorm:"size:1"` // A, B, C, D, E, or NULL
	IsCorrect      *bool      `json:"isCorrect"`
	AnsweredAt     *time.Time `json:"answeredAt"`

	// Relations
	Attempt  TryOutAttempt  `json:"attempt,omitempty" gorm:"foreignKey:AttemptID"`
	Question TryOutQuestion `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}
