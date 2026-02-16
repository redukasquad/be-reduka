package answers

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type answerService struct {
	repo Repository
}

type Service interface {
	GetByRegistrationID(registrationID uint, requestID string) ([]AnswerResponse, error)
	CreateAnswer(input CreateAnswerRequest, requestID string) (*AnswerResponse, error)
	DeleteAnswer(id uint, requestID string) error
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

func (s *answerService) CreateAnswer(input CreateAnswerRequest, requestID string) (*AnswerResponse, error) {
	utils.LogInfo("answers", "create", "Creating answer", requestID, 0, map[string]any{
		"registration_id": input.RegistrationID,
		"question_id":     input.QuestionID,
	})

	answer := &entities.RegistrationAnswer{
		RegistrationID: input.RegistrationID,
		QuestionID:     input.QuestionID,
		AnswerText:     input.AnswerText,
	}

	if err := s.repo.Create(answer); err != nil {
		utils.LogError("answers", "create", "Failed to create answer: "+err.Error(), requestID, 0, map[string]any{
			"registration_id": input.RegistrationID,
			"question_id":     input.QuestionID,
		})
		return nil, err
	}

	created, err := s.repo.FindByID(answer.ID)
	if err != nil {
		utils.LogError("answers", "create", "Failed to fetch created answer: "+err.Error(), requestID, 0, map[string]any{
			"answer_id": answer.ID,
		})
		return nil, err
	}

	response := s.toAnswerResponse(created)

	utils.LogSuccess("answers", "create", "Answer created successfully", requestID, 0, map[string]any{
		"answer_id":       answer.ID,
		"registration_id": input.RegistrationID,
		"question_id":     input.QuestionID,
	})

	return &response, nil
}

func (s *answerService) DeleteAnswer(id uint, requestID string) error {
	utils.LogInfo("answers", "delete", "Deleting answer", requestID, 0, map[string]any{
		"answer_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("answers", "delete", "Answer not found", requestID, 0, map[string]any{
				"answer_id": id,
			})
			return errors.New("answer not found")
		}
		utils.LogError("answers", "delete", "Failed to find answer: "+err.Error(), requestID, 0, map[string]any{
			"answer_id": id,
		})
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("answers", "delete", "Failed to delete answer: "+err.Error(), requestID, 0, map[string]any{
			"answer_id": id,
		})
		return err
	}

	utils.LogSuccess("answers", "delete", "Answer deleted successfully", requestID, 0, map[string]any{
		"answer_id": id,
	})

	return nil
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
