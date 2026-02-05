package entities

import (
	"time"

	"gorm.io/gorm"
)

// AttemptStatus represents the status of a Try Out attempt.
type AttemptStatus string

const (
	AttemptStatusNotStarted AttemptStatus = "not_started"
	AttemptStatusInProgress AttemptStatus = "in_progress"
	AttemptStatusCompleted  AttemptStatus = "completed"
)

// TryOutAttempt represents a user's attempt session for a Try Out.
// Each registration can only have one attempt (user can only take each Try Out once).
type TryOutAttempt struct {
	gorm.Model

	RegistrationID uint `json:"registrationId" gorm:"uniqueIndex;not null"` // 1:1 with registration

	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`

	CurrentSubtestID *uint         `json:"currentSubtestId"` // Tracking progress
	Status           AttemptStatus `json:"status" gorm:"size:20;default:'not_started'"`
	TotalScore       *float64      `json:"totalScore" gorm:"type:decimal(10,2)"` // For leaderboard

	// Relations
	Registration   TryOutRegistration `json:"registration,omitempty" gorm:"foreignKey:RegistrationID"`
	CurrentSubtest *Subtest           `json:"currentSubtest,omitempty" gorm:"foreignKey:CurrentSubtestID"`
	SubtestResults []SubtestResult    `json:"subtestResults,omitempty" gorm:"foreignKey:AttemptID"`
	Answers        []UserTryOutAnswer `json:"answers,omitempty" gorm:"foreignKey:AttemptID"`
}
