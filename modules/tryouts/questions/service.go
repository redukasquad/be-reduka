package questions

import (
	"errors"
	"fmt"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type questionService struct {
	repo       Repository
	tryOutRepo TryOutRepository
}

// TryOutRepository is an interface for try out related operations
type TryOutRepository interface {
	HasTutorPermission(tryOutID, userID uint) (bool, error)
	FindByID(id uint) (entities.TryOut, error)
}

type Service interface {
	// Subtests
	GetSubtestsWithQuestionCount(tryOutID uint, requestID string) ([]SubtestWithQuestionsResponse, error)

	// Questions
	GetQuestionsByTryOut(tryOutID uint, requestID string) ([]QuestionResponse, error)
	GetQuestionsBySubtest(tryOutID, subtestID uint, requestID string) ([]QuestionResponse, error)
	GetQuestionByID(id uint, requestID string) (*QuestionResponse, error)
	CreateQuestion(tryOutID, subtestID uint, input CreateQuestionInput, requestID string, userID uint) (*QuestionResponse, error)
	UpdateQuestion(id uint, input UpdateQuestionInput, requestID string, userID uint) (*QuestionResponse, error)
	DeleteQuestion(id uint, requestID string, userID uint) error

	// Permission check
	CheckTutorPermission(tryOutID, userID uint) error
}

func NewService(repo Repository, tryOutRepo TryOutRepository) Service {
	return &questionService{
		repo:       repo,
		tryOutRepo: tryOutRepo,
	}
}

// ==========================================
// Permission Check
// ==========================================

func (s *questionService) CheckTutorPermission(tryOutID, userID uint) error {
	hasPermission, err := s.tryOutRepo.HasTutorPermission(tryOutID, userID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("you don't have permission to manage questions for this try out")
	}
	return nil
}

// ==========================================
// Subtest Service Methods
// ==========================================

func (s *questionService) GetSubtestsWithQuestionCount(tryOutID uint, requestID string) ([]SubtestWithQuestionsResponse, error) {
	utils.LogInfo("questions", "get_subtests", "Fetching subtests with question count", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
	})

	subtests, err := s.repo.FindAllSubtests()
	if err != nil {
		utils.LogError("questions", "get_subtests", "Failed to fetch subtests: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []SubtestWithQuestionsResponse
	for _, subtest := range subtests {
		count, err := s.repo.CountByTryOutAndSubtest(tryOutID, subtest.ID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, SubtestWithQuestionsResponse{
			SubtestBriefResponse: ToSubtestBriefResponse(subtest),
			CurrentQuestionCount: int(count),
			IsComplete:           int(count) >= subtest.QuestionCount,
		})
	}

	return responses, nil
}

// ==========================================
// Question Service Methods
// ==========================================

func (s *questionService) GetQuestionsByTryOut(tryOutID uint, requestID string) ([]QuestionResponse, error) {
	utils.LogInfo("questions", "get_by_tryout", "Fetching questions by try out", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
	})

	questions, err := s.repo.FindByTryOutID(tryOutID)
	if err != nil {
		utils.LogError("questions", "get_by_tryout", "Failed to fetch questions: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []QuestionResponse
	for _, q := range questions {
		responses = append(responses, ToQuestionResponse(q))
	}

	return responses, nil
}

func (s *questionService) GetQuestionsBySubtest(tryOutID, subtestID uint, requestID string) ([]QuestionResponse, error) {
	utils.LogInfo("questions", "get_by_subtest", "Fetching questions by subtest", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
		"subtest_id": subtestID,
	})

	questions, err := s.repo.FindByTryOutAndSubtest(tryOutID, subtestID)
	if err != nil {
		utils.LogError("questions", "get_by_subtest", "Failed to fetch questions: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []QuestionResponse
	for _, q := range questions {
		responses = append(responses, ToQuestionResponse(q))
	}

	return responses, nil
}

func (s *questionService) GetQuestionByID(id uint, requestID string) (*QuestionResponse, error) {
	utils.LogInfo("questions", "get_by_id", "Fetching question by ID", requestID, 0, map[string]any{
		"question_id": id,
	})

	question, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		utils.LogError("questions", "get_by_id", "Failed to fetch question: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	response := ToQuestionResponse(question)
	return &response, nil
}

func (s *questionService) CreateQuestion(tryOutID, subtestID uint, input CreateQuestionInput, requestID string, userID uint) (*QuestionResponse, error) {
	utils.LogInfo("questions", "create", "Attempting to create question", requestID, userID, map[string]any{
		"try_out_id": tryOutID,
		"subtest_id": subtestID,
	})

	// Check try out exists
	_, err := s.tryOutRepo.FindByID(tryOutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("try out not found")
		}
		return nil, err
	}

	// Check subtest exists and get question limit
	subtest, err := s.repo.FindSubtestByID(subtestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subtest not found")
		}
		return nil, err
	}

	// Check question count limit
	currentCount, err := s.repo.CountByTryOutAndSubtest(tryOutID, subtestID)
	if err != nil {
		return nil, err
	}

	if int(currentCount) >= subtest.QuestionCount {
		return nil, fmt.Errorf("subtest %s already has maximum %d questions", subtest.Code, subtest.QuestionCount)
	}

	// Validate order number
	if input.OrderNumber > subtest.QuestionCount {
		return nil, fmt.Errorf("order number cannot exceed %d for subtest %s", subtest.QuestionCount, subtest.Code)
	}

	question := &entities.TryOutQuestion{
		TryOutPackageID: tryOutID,
		SubtestID:       subtestID,
		QuestionText:    input.QuestionText,
		ImageURL:        input.ImageURL,
		Explanation:     input.Explanation,
		DifficultyLevel: entities.DifficultyLevel(input.DifficultyLevel),
		OrderNumber:     input.OrderNumber,
		OptionA:         input.OptionA,
		OptionB:         input.OptionB,
		OptionC:         input.OptionC,
		OptionD:         input.OptionD,
		OptionE:         input.OptionE,
		CorrectOption:   input.CorrectOption,
		CreatedByUserID: userID,
	}

	if err := s.repo.Create(question); err != nil {
		utils.LogError("questions", "create", "Failed to create question: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	// Fetch with preload
	createdQuestion, err := s.repo.FindByID(question.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("questions", "create", "Question created successfully", requestID, userID, map[string]any{
		"question_id": question.ID,
	})

	response := ToQuestionResponse(createdQuestion)
	return &response, nil
}

func (s *questionService) UpdateQuestion(id uint, input UpdateQuestionInput, requestID string, userID uint) (*QuestionResponse, error) {
	utils.LogInfo("questions", "update", "Attempting to update question", requestID, userID, map[string]any{
		"question_id": id,
	})

	question, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}

	// Update fields if provided
	if input.QuestionText != nil {
		question.QuestionText = *input.QuestionText
	}
	if input.ImageURL != nil {
		question.ImageURL = *input.ImageURL
	}
	if input.Explanation != nil {
		question.Explanation = *input.Explanation
	}
	if input.DifficultyLevel != nil {
		question.DifficultyLevel = entities.DifficultyLevel(*input.DifficultyLevel)
	}
	if input.OrderNumber != nil {
		question.OrderNumber = *input.OrderNumber
	}
	if input.OptionA != nil {
		question.OptionA = *input.OptionA
	}
	if input.OptionB != nil {
		question.OptionB = *input.OptionB
	}
	if input.OptionC != nil {
		question.OptionC = *input.OptionC
	}
	if input.OptionD != nil {
		question.OptionD = *input.OptionD
	}
	if input.OptionE != nil {
		question.OptionE = *input.OptionE
	}
	if input.CorrectOption != nil {
		question.CorrectOption = *input.CorrectOption
	}

	if err := s.repo.Update(&question); err != nil {
		utils.LogError("questions", "update", "Failed to update question: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	// Fetch with preload
	updatedQuestion, err := s.repo.FindByID(question.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("questions", "update", "Question updated successfully", requestID, userID, map[string]any{
		"question_id": question.ID,
	})

	response := ToQuestionResponse(updatedQuestion)
	return &response, nil
}

func (s *questionService) DeleteQuestion(id uint, requestID string, userID uint) error {
	utils.LogInfo("questions", "delete", "Attempting to delete question", requestID, userID, map[string]any{
		"question_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("question not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("questions", "delete", "Failed to delete question: "+err.Error(), requestID, userID, nil)
		return err
	}

	utils.LogSuccess("questions", "delete", "Question deleted successfully", requestID, userID, map[string]any{
		"question_id": id,
	})
	return nil
}
