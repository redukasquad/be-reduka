package answers

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
)

type answerService struct {
	repo Repository
}

type Service interface {
	GetByRegistrationID(registrationID uint, requestID string) ([]AnswerResponse, error)
}

func NewService(repo Repository) Service {
	return &answerService{repo: repo}
}

func (s *answerService) GetByRegistrationID(registrationID uint, requestID string) ([]AnswerResponse, error) {
	utils.LogInfo("answers", "get_by_registration", "Fetching answers for registration", requestID, 0, map[string]any{
		"registration_id": registrationID,
	})

	answers, err := s.repo.FindByRegistrationID(registrationID)
	if err != nil {
		utils.LogError("answers", "get_by_registration", "Failed to fetch answers: "+err.Error(), requestID, 0, map[string]any{
			"registration_id": registrationID,
		})
		return nil, err
	}

	var responses []AnswerResponse
	for _, a := range answers {
		responses = append(responses, s.toAnswerResponse(a))
	}

	utils.LogSuccess("answers", "get_by_registration", "Successfully fetched answers", requestID, 0, map[string]any{
		"registration_id": registrationID,
		"count":           len(responses),
	})
	return responses, nil
}

func (s *answerService) toAnswerResponse(a entities.RegistrationAnswer) AnswerResponse {
	return AnswerResponse{
		ID:             a.ID,
		RegistrationID: a.RegistrationID,
		QuestionID:     a.QuestionID,
		QuestionText:   a.Question.QuestionText,
		AnswerText:     a.AnswerText,
	}
}
