package questions

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type questionService struct {
	repo Repository
}

// Service interface defines the business logic methods for questions
type Service interface {
	GetByCourseID(courseID uint, requestID string) ([]QuestionResponse, error)
	Create(courseID uint, input CreateQuestionInput, requestID string, userID uint) (*QuestionResponse, error)
	Update(id uint, input UpdateQuestionInput, requestID string, userID uint) (*QuestionResponse, error)
	Delete(id uint, requestID string, userID uint) error
}

// NewService creates a new question service
func NewService(repo Repository) Service {
	return &questionService{repo: repo}
}

func (s *questionService) GetByCourseID(courseID uint, requestID string) ([]QuestionResponse, error) {
	utils.LogInfo("questions", "get_by_course", "Fetching questions for course", requestID, 0, map[string]any{
		"course_id": courseID,
	})

	questions, err := s.repo.FindByCourseID(courseID)
	if err != nil {
		utils.LogError("questions", "get_by_course", "Failed to fetch questions: "+err.Error(), requestID, 0, map[string]any{
			"course_id": courseID,
		})
		return nil, err
	}

	var responses []QuestionResponse
	for _, q := range questions {
		responses = append(responses, s.toQuestionResponse(q))
	}

	utils.LogSuccess("questions", "get_by_course", "Successfully fetched questions", requestID, 0, map[string]any{
		"course_id": courseID,
		"count":     len(responses),
	})
	return responses, nil
}

func (s *questionService) Create(courseID uint, input CreateQuestionInput, requestID string, userID uint) (*QuestionResponse, error) {
	utils.LogInfo("questions", "create", "Creating new question", requestID, userID, map[string]any{
		"course_id": courseID,
	})

	question := &entities.RegistrationQuestion{
		CourseID:      courseID,
		QuestionText:  input.QuestionText,
		QuestionType:  input.QuestionType,
		QuestionOrder: input.QuestionOrder,
	}

	if err := s.repo.Create(question); err != nil {
		utils.LogError("questions", "create", "Failed to create question: "+err.Error(), requestID, userID, map[string]any{
			"course_id": courseID,
		})
		return nil, err
	}

	utils.LogSuccess("questions", "create", "Question created successfully", requestID, userID, map[string]any{
		"course_id":   courseID,
		"question_id": question.ID,
	})

	response := s.toQuestionResponse(*question)
	return &response, nil
}

func (s *questionService) Update(id uint, input UpdateQuestionInput, requestID string, userID uint) (*QuestionResponse, error) {
	utils.LogInfo("questions", "update", "Updating question", requestID, userID, map[string]any{
		"question_id": id,
	})

	question, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("questions", "update", "Question not found", requestID, userID, map[string]any{
				"question_id": id,
			})
			return nil, errors.New("question not found")
		}
		return nil, err
	}

	if input.QuestionText != nil {
		question.QuestionText = *input.QuestionText
	}
	if input.QuestionType != nil {
		question.QuestionType = *input.QuestionType
	}
	if input.QuestionOrder != nil {
		question.QuestionOrder = *input.QuestionOrder
	}

	if err := s.repo.Update(&question); err != nil {
		utils.LogError("questions", "update", "Failed to update question: "+err.Error(), requestID, userID, map[string]any{
			"question_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("questions", "update", "Question updated successfully", requestID, userID, map[string]any{
		"question_id": id,
	})

	response := s.toQuestionResponse(question)
	return &response, nil
}

func (s *questionService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("questions", "delete", "Deleting question", requestID, userID, map[string]any{
		"question_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("questions", "delete", "Question not found", requestID, userID, map[string]any{
				"question_id": id,
			})
			return errors.New("question not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("questions", "delete", "Failed to delete question: "+err.Error(), requestID, userID, map[string]any{
			"question_id": id,
		})
		return err
	}

	utils.LogSuccess("questions", "delete", "Question deleted successfully", requestID, userID, map[string]any{
		"question_id": id,
	})
	return nil
}

func (s *questionService) toQuestionResponse(q entities.RegistrationQuestion) QuestionResponse {
	return QuestionResponse{
		ID:            q.ID,
		CourseID:      q.CourseID,
		QuestionText:  q.QuestionText,
		QuestionType:  q.QuestionType,
		QuestionOrder: q.QuestionOrder,
	}
}
