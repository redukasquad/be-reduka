package attempts

import (
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
)

// ==========================================
// ATTEMPT DTOs
// ==========================================

// AttemptResponse is the full response for an attempt
type AttemptResponse struct {
	ID               uint                    `json:"id"`
	RegistrationID   uint                    `json:"registrationId"`
	TryOut           *TryOutBriefResponse    `json:"tryOut,omitempty"`
	StartedAt        *time.Time              `json:"startedAt,omitempty"`
	FinishedAt       *time.Time              `json:"finishedAt,omitempty"`
	Status           string                  `json:"status"`
	CurrentSubtestID *uint                   `json:"currentSubtestId,omitempty"`
	TotalScore       *float64                `json:"totalScore,omitempty"`
	SubtestResults   []SubtestResultResponse `json:"subtestResults,omitempty"`
}

// AttemptCurrentStateResponse shows current state during exam
type AttemptCurrentStateResponse struct {
	ID             uint                      `json:"id"`
	Status         string                    `json:"status"`
	CurrentSubtest *SubtestBriefResponse     `json:"currentSubtest,omitempty"`
	TimeRemaining  *int                      `json:"timeRemaining,omitempty"` // in seconds
	SubtestResults []SubtestProgressResponse `json:"subtestProgress,omitempty"`
}

// SubtestResultResponse shows subtest result after completion
type SubtestResultResponse struct {
	ID              uint                  `json:"id"`
	SubtestID       uint                  `json:"subtestId"`
	Subtest         *SubtestBriefResponse `json:"subtest,omitempty"`
	StartedAt       *time.Time            `json:"startedAt,omitempty"`
	FinishedAt      *time.Time            `json:"finishedAt,omitempty"`
	CorrectCount    int                   `json:"correctCount"`
	WrongCount      int                   `json:"wrongCount"`
	UnansweredCount int                   `json:"unansweredCount"`
	RawScore        *float64              `json:"rawScore,omitempty"`
	FinalScore      *float64              `json:"finalScore,omitempty"`
}

// SubtestProgressResponse shows progress during exam
type SubtestProgressResponse struct {
	SubtestID     uint   `json:"subtestId"`
	SubtestCode   string `json:"subtestCode"`
	SubtestName   string `json:"subtestName"`
	Status        string `json:"status"` // not_started, in_progress, completed
	AnsweredCount int    `json:"answeredCount"`
	TotalCount    int    `json:"totalCount"`
}

// SubtestBriefResponse is a minimal subtest info
type SubtestBriefResponse struct {
	ID               uint    `json:"id"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	QuestionCount    int     `json:"questionCount"`
	TimeLimitSeconds int     `json:"timeLimitSeconds"`
	MaxScore         float64 `json:"maxScore"`
}

// TryOutBriefResponse is a minimal try out info
type TryOutBriefResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// QuestionForExamResponse shows question during exam (without correct answer)
type QuestionForExamResponse struct {
	ID              uint    `json:"id"`
	OrderNumber     int     `json:"orderNumber"`
	QuestionText    string  `json:"questionText"`
	ImageURL        string  `json:"imageUrl,omitempty"`
	DifficultyLevel string  `json:"difficultyLevel"`
	OptionA         string  `json:"optionA"`
	OptionB         string  `json:"optionB"`
	OptionC         string  `json:"optionC"`
	OptionD         string  `json:"optionD"`
	OptionE         string  `json:"optionE"`
	SelectedOption  *string `json:"selectedOption,omitempty"` // User's current answer
}

// LeaderboardEntryResponse shows leaderboard entry
type LeaderboardEntryResponse struct {
	Rank       int        `json:"rank"`
	UserID     uint       `json:"userId"`
	Username   string     `json:"username"`
	TotalScore float64    `json:"totalScore"`
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
}

// SubmitAnswerInput is the input for submitting an answer
type SubmitAnswerInput struct {
	QuestionID     uint   `json:"questionId" binding:"required"`
	SelectedOption string `json:"selectedOption" binding:"omitempty,oneof=A B C D E"`
}

// SubmitSubtestInput is the input for submitting all answers for a subtest
type SubmitSubtestInput struct {
	Answers []SubmitAnswerInput `json:"answers" binding:"required,dive"`
}

// ==========================================
// Helper Functions
// ==========================================

func ToAttemptResponse(a entities.TryOutAttempt) AttemptResponse {
	response := AttemptResponse{
		ID:               a.ID,
		RegistrationID:   a.RegistrationID,
		StartedAt:        a.StartedAt,
		FinishedAt:       a.FinishedAt,
		Status:           string(a.Status),
		CurrentSubtestID: a.CurrentSubtestID,
		TotalScore:       a.TotalScore,
	}

	if a.Registration.TryOutPackage.ID != 0 {
		tryOut := ToTryOutBriefResponse(a.Registration.TryOutPackage)
		response.TryOut = &tryOut
	}

	for _, sr := range a.SubtestResults {
		response.SubtestResults = append(response.SubtestResults, ToSubtestResultResponse(sr))
	}

	return response
}

func ToSubtestResultResponse(sr entities.SubtestResult) SubtestResultResponse {
	response := SubtestResultResponse{
		ID:              sr.ID,
		SubtestID:       sr.SubtestID,
		StartedAt:       sr.StartedAt,
		FinishedAt:      sr.FinishedAt,
		CorrectCount:    sr.CorrectCount,
		WrongCount:      sr.WrongCount,
		UnansweredCount: sr.UnansweredCount,
		RawScore:        sr.RawScore,
		FinalScore:      sr.FinalScore,
	}

	if sr.Subtest.ID != 0 {
		subtest := ToSubtestBriefResponse(sr.Subtest)
		response.Subtest = &subtest
	}

	return response
}

func ToSubtestBriefResponse(s entities.Subtest) SubtestBriefResponse {
	return SubtestBriefResponse{
		ID:               s.ID,
		Code:             s.Code,
		Name:             s.Name,
		QuestionCount:    s.QuestionCount,
		TimeLimitSeconds: s.TimeLimitSeconds,
		MaxScore:         s.MaxScore,
	}
}

func ToTryOutBriefResponse(t entities.TryOut) TryOutBriefResponse {
	return TryOutBriefResponse{
		ID:   t.ID,
		Name: t.Name,
	}
}

func ToQuestionForExamResponse(q entities.TryOutQuestion, selectedOption *string) QuestionForExamResponse {
	return QuestionForExamResponse{
		ID:              q.ID,
		OrderNumber:     q.OrderNumber,
		QuestionText:    q.QuestionText,
		ImageURL:        q.ImageURL,
		DifficultyLevel: string(q.DifficultyLevel),
		OptionA:         q.OptionA,
		OptionB:         q.OptionB,
		OptionC:         q.OptionC,
		OptionD:         q.OptionD,
		OptionE:         q.OptionE,
		SelectedOption:  selectedOption,
	}
}
