package questions

import (
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
)

// ==========================================
// QUESTION DTOs
// ==========================================

// QuestionResponse is the response DTO for a question
type QuestionResponse struct {
	ID              uint                  `json:"id"`
	TryOutID        uint                  `json:"tryOutId"`
	SubtestID       uint                  `json:"subtestId"`
	Subtest         *SubtestBriefResponse `json:"subtest,omitempty"`
	QuestionText    string                `json:"questionText"`
	ImageURL        string                `json:"imageUrl,omitempty"`
	Explanation     string                `json:"explanation,omitempty"`
	DifficultyLevel string                `json:"difficultyLevel"`
	OrderNumber     int                   `json:"orderNumber"`
	OptionA         string                `json:"optionA"`
	OptionB         string                `json:"optionB"`
	OptionC         string                `json:"optionC"`
	OptionD         string                `json:"optionD"`
	OptionE         string                `json:"optionE"`
	CorrectOption   string                `json:"correctOption"`
	CreatedAt       time.Time             `json:"createdAt"`
}

// QuestionBriefResponse is a minimal response for list views (without correct answer for students)
type QuestionBriefResponse struct {
	ID              uint   `json:"id"`
	OrderNumber     int    `json:"orderNumber"`
	QuestionText    string `json:"questionText"`
	ImageURL        string `json:"imageUrl,omitempty"`
	DifficultyLevel string `json:"difficultyLevel"`
	OptionA         string `json:"optionA"`
	OptionB         string `json:"optionB"`
	OptionC         string `json:"optionC"`
	OptionD         string `json:"optionD"`
	OptionE         string `json:"optionE"`
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

// SubtestWithQuestionsResponse shows subtest with question count
type SubtestWithQuestionsResponse struct {
	SubtestBriefResponse
	CurrentQuestionCount int  `json:"currentQuestionCount"`
	IsComplete           bool `json:"isComplete"`
}

// CreateQuestionInput is the input for creating a new question
type CreateQuestionInput struct {
	QuestionText    string `json:"questionText" binding:"required"`
	ImageURL        string `json:"imageUrl"`
	Explanation     string `json:"explanation"`
	DifficultyLevel string `json:"difficultyLevel" binding:"required,oneof=easy medium hard"`
	OrderNumber     int    `json:"orderNumber" binding:"required,min=1"`
	OptionA         string `json:"optionA" binding:"required"`
	OptionB         string `json:"optionB" binding:"required"`
	OptionC         string `json:"optionC" binding:"required"`
	OptionD         string `json:"optionD" binding:"required"`
	OptionE         string `json:"optionE" binding:"required"`
	CorrectOption   string `json:"correctOption" binding:"required,oneof=A B C D E"`
}

// UpdateQuestionInput is the input for updating a question
type UpdateQuestionInput struct {
	QuestionText    *string `json:"questionText"`
	ImageURL        *string `json:"imageUrl"`
	Explanation     *string `json:"explanation"`
	DifficultyLevel *string `json:"difficultyLevel" binding:"omitempty,oneof=easy medium hard"`
	OrderNumber     *int    `json:"orderNumber" binding:"omitempty,min=1"`
	OptionA         *string `json:"optionA"`
	OptionB         *string `json:"optionB"`
	OptionC         *string `json:"optionC"`
	OptionD         *string `json:"optionD"`
	OptionE         *string `json:"optionE"`
	CorrectOption   *string `json:"correctOption" binding:"omitempty,oneof=A B C D E"`
}

// ==========================================
// Helper Functions
// ==========================================

func ToQuestionResponse(q entities.TryOutQuestion) QuestionResponse {
	response := QuestionResponse{
		ID:              q.ID,
		TryOutID:        q.TryOutPackageID,
		SubtestID:       q.SubtestID,
		QuestionText:    q.QuestionText,
		ImageURL:        q.ImageURL,
		Explanation:     q.Explanation,
		DifficultyLevel: string(q.DifficultyLevel),
		OrderNumber:     q.OrderNumber,
		OptionA:         q.OptionA,
		OptionB:         q.OptionB,
		OptionC:         q.OptionC,
		OptionD:         q.OptionD,
		OptionE:         q.OptionE,
		CorrectOption:   q.CorrectOption,
		CreatedAt:       q.CreatedAt,
	}

	if q.Subtest.ID != 0 {
		subtest := ToSubtestBriefResponse(q.Subtest)
		response.Subtest = &subtest
	}

	return response
}

func ToQuestionBriefResponse(q entities.TryOutQuestion) QuestionBriefResponse {
	return QuestionBriefResponse{
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
	}
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
