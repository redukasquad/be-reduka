package attempts

import (
	"errors"
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

// Difficulty weights for scoring
const (
	WeightEasy   = 1.0
	WeightMedium = 1.5
	WeightHard   = 2.0
)

type attemptService struct {
	repo Repository
}

type Service interface {
	// Attempt lifecycle
	StartAttempt(registrationID uint, userID uint, requestID string) (*AttemptResponse, error)
	GetCurrentState(attemptID uint, userID uint, requestID string) (*AttemptCurrentStateResponse, error)
	GetAttemptResults(attemptID uint, userID uint, requestID string) (*AttemptResponse, error)

	// Subtest operations
	StartSubtest(attemptID, subtestID uint, userID uint, requestID string) ([]QuestionForExamResponse, error)
	SubmitSubtest(attemptID, subtestID uint, input SubmitSubtestInput, userID uint, requestID string) (*SubtestResultResponse, error)

	// Finish attempt
	FinishAttempt(attemptID uint, userID uint, requestID string) (*AttemptResponse, error)

	// Leaderboard
	GetLeaderboard(tryOutID uint, requestID string) ([]LeaderboardEntryResponse, error)
}

func NewService(repo Repository) Service {
	return &attemptService{repo: repo}
}

// ==========================================
// Attempt Lifecycle
// ==========================================

func (s *attemptService) StartAttempt(registrationID uint, userID uint, requestID string) (*AttemptResponse, error) {
	utils.LogInfo("attempts", "start", "Starting attempt", requestID, userID, map[string]any{
		"registration_id": registrationID,
	})

	// Get registration
	registration, err := s.repo.FindRegistrationByID(registrationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	// Check ownership
	if registration.UserID != userID {
		return nil, errors.New("you can only start your own registration")
	}

	// Check if payment is approved
	if registration.PaymentStatus != entities.PaymentStatusApproved {
		return nil, errors.New("payment must be approved before starting")
	}

	// Check if attempt already exists
	existingAttempt, err := s.repo.FindAttemptByRegistrationID(registrationID)
	if err == nil && existingAttempt.ID != 0 {
		// Return existing attempt
		response := ToAttemptResponse(existingAttempt)
		return &response, nil
	}

	// Create new attempt
	now := time.Now()
	subtests, err := s.repo.FindAllSubtests()
	if err != nil {
		return nil, err
	}

	var firstSubtestID *uint
	if len(subtests) > 0 {
		firstSubtestID = &subtests[0].ID
	}

	attempt := &entities.TryOutAttempt{
		RegistrationID:   registrationID,
		StartedAt:        &now,
		Status:           entities.AttemptStatusInProgress,
		CurrentSubtestID: firstSubtestID,
	}

	if err := s.repo.CreateAttempt(attempt); err != nil {
		utils.LogError("attempts", "start", "Failed to create attempt: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	// Fetch with preload
	createdAttempt, err := s.repo.FindAttemptByID(attempt.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("attempts", "start", "Attempt started", requestID, userID, map[string]any{
		"attempt_id": attempt.ID,
	})

	response := ToAttemptResponse(createdAttempt)
	return &response, nil
}

func (s *attemptService) GetCurrentState(attemptID uint, userID uint, requestID string) (*AttemptCurrentStateResponse, error) {
	attempt, err := s.repo.FindAttemptByID(attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	// Check ownership
	if attempt.Registration.UserID != userID {
		return nil, errors.New("you can only view your own attempt")
	}

	response := &AttemptCurrentStateResponse{
		ID:     attempt.ID,
		Status: string(attempt.Status),
	}

	// Get current subtest info
	if attempt.CurrentSubtest != nil {
		subtest := ToSubtestBriefResponse(*attempt.CurrentSubtest)
		response.CurrentSubtest = &subtest

		// Calculate time remaining for current subtest
		currentResult, err := s.repo.FindSubtestResultByAttemptAndSubtest(attemptID, attempt.CurrentSubtest.ID)
		if err == nil && currentResult.StartedAt != nil {
			elapsed := int(time.Since(*currentResult.StartedAt).Seconds())
			remaining := attempt.CurrentSubtest.TimeLimitSeconds - elapsed
			if remaining < 0 {
				remaining = 0
			}
			response.TimeRemaining = &remaining
		}
	}

	// Get progress for all subtests
	subtests, _ := s.repo.FindAllSubtests()
	for _, subtest := range subtests {
		progress := SubtestProgressResponse{
			SubtestID:   subtest.ID,
			SubtestCode: subtest.Code,
			SubtestName: subtest.Name,
			Status:      "not_started",
			TotalCount:  subtest.QuestionCount,
		}

		result, err := s.repo.FindSubtestResultByAttemptAndSubtest(attemptID, subtest.ID)
		if err == nil {
			if result.FinishedAt != nil {
				progress.Status = "completed"
				progress.AnsweredCount = result.CorrectCount + result.WrongCount
			} else if result.StartedAt != nil {
				progress.Status = "in_progress"
				// Count answered
				answers, _ := s.repo.FindAnswersByAttemptAndSubtest(attemptID, subtest.ID)
				for _, ans := range answers {
					if ans.SelectedOption != nil {
						progress.AnsweredCount++
					}
				}
			}
		}

		response.SubtestResults = append(response.SubtestResults, progress)
	}

	return response, nil
}

func (s *attemptService) GetAttemptResults(attemptID uint, userID uint, requestID string) (*AttemptResponse, error) {
	attempt, err := s.repo.FindAttemptByID(attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	// Check ownership
	if attempt.Registration.UserID != userID {
		return nil, errors.New("you can only view your own results")
	}

	if attempt.Status != entities.AttemptStatusCompleted {
		return nil, errors.New("attempt is not completed yet")
	}

	response := ToAttemptResponse(attempt)
	return &response, nil
}

// ==========================================
// Subtest Operations
// ==========================================

func (s *attemptService) StartSubtest(attemptID, subtestID uint, userID uint, requestID string) ([]QuestionForExamResponse, error) {
	utils.LogInfo("attempts", "start_subtest", "Starting subtest", requestID, userID, map[string]any{
		"attempt_id": attemptID,
		"subtest_id": subtestID,
	})

	attempt, err := s.repo.FindAttemptByID(attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	// Check ownership
	if attempt.Registration.UserID != userID {
		return nil, errors.New("you can only access your own attempt")
	}

	// Check attempt status
	if attempt.Status != entities.AttemptStatusInProgress {
		return nil, errors.New("attempt is not in progress")
	}

	// Check subtest exists
	subtest, err := s.repo.FindSubtestByID(subtestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subtest not found")
		}
		return nil, err
	}

	// Check if subtest result exists, create if not
	result, err := s.repo.FindSubtestResultByAttemptAndSubtest(attemptID, subtestID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		now := time.Now()
		result = entities.SubtestResult{
			AttemptID: attemptID,
			SubtestID: subtestID,
			StartedAt: &now,
		}
		if err := s.repo.CreateSubtestResult(&result); err != nil {
			return nil, err
		}
	}

	// Update current subtest
	attempt.CurrentSubtestID = &subtestID
	if err := s.repo.UpdateAttempt(&attempt); err != nil {
		return nil, err
	}

	// Get questions for this subtest
	tryOutID := attempt.Registration.TryOutPackageID
	questions, err := s.repo.FindQuestionsByTryOutAndSubtest(tryOutID, subtestID)
	if err != nil {
		return nil, err
	}

	// Get existing answers
	answerMap := make(map[uint]*string)
	answers, _ := s.repo.FindAnswersByAttemptAndSubtest(attemptID, subtestID)
	for _, ans := range answers {
		answerMap[ans.QuestionID] = ans.SelectedOption
	}

	var responses []QuestionForExamResponse
	for _, q := range questions {
		responses = append(responses, ToQuestionForExamResponse(q, answerMap[q.ID]))
	}

	utils.LogSuccess("attempts", "start_subtest", "Subtest started", requestID, userID, map[string]any{
		"attempt_id":     attemptID,
		"subtest_id":     subtestID,
		"subtest_code":   subtest.Code,
		"question_count": len(questions),
	})

	return responses, nil
}

func (s *attemptService) SubmitSubtest(attemptID, subtestID uint, input SubmitSubtestInput, userID uint, requestID string) (*SubtestResultResponse, error) {
	utils.LogInfo("attempts", "submit_subtest", "Submitting subtest", requestID, userID, map[string]any{
		"attempt_id":   attemptID,
		"subtest_id":   subtestID,
		"answer_count": len(input.Answers),
	})

	attempt, err := s.repo.FindAttemptByID(attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	// Check ownership
	if attempt.Registration.UserID != userID {
		return nil, errors.New("you can only submit your own answers")
	}

	// Check attempt status
	if attempt.Status != entities.AttemptStatusInProgress {
		return nil, errors.New("attempt is not in progress")
	}

	// Get subtest result
	result, err := s.repo.FindSubtestResultByAttemptAndSubtest(attemptID, subtestID)
	if err != nil {
		return nil, errors.New("subtest not started")
	}

	if result.FinishedAt != nil {
		return nil, errors.New("subtest already submitted")
	}

	// Get subtest for scoring
	subtest, err := s.repo.FindSubtestByID(subtestID)
	if err != nil {
		return nil, err
	}

	// Process answers
	var correctCount, wrongCount, unansweredCount int
	var rawScore float64

	tryOutID := attempt.Registration.TryOutPackageID
	questions, err := s.repo.FindQuestionsByTryOutAndSubtest(tryOutID, subtestID)
	if err != nil {
		return nil, err
	}

	questionMap := make(map[uint]entities.TryOutQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	answeredQuestions := make(map[uint]bool)

	// Save each answer
	for _, ans := range input.Answers {
		question, exists := questionMap[ans.QuestionID]
		if !exists {
			continue // Skip invalid question
		}

		answeredQuestions[ans.QuestionID] = true

		var isCorrect *bool
		now := time.Now()

		if ans.SelectedOption != "" {
			correct := ans.SelectedOption == question.CorrectOption
			isCorrect = &correct

			if correct {
				correctCount++
				// Calculate weighted score
				switch question.DifficultyLevel {
				case entities.DifficultyEasy:
					rawScore += WeightEasy
				case entities.DifficultyMedium:
					rawScore += WeightMedium
				case entities.DifficultyHard:
					rawScore += WeightHard
				}
			} else {
				wrongCount++
			}
		}

		// Upsert answer
		selectedOption := &ans.SelectedOption
		if ans.SelectedOption == "" {
			selectedOption = nil
		}

		existingAnswer, err := s.repo.FindAnswerByAttemptAndQuestion(attemptID, ans.QuestionID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			answer := entities.UserTryOutAnswer{
				AttemptID:      attemptID,
				QuestionID:     ans.QuestionID,
				SelectedOption: selectedOption,
				IsCorrect:      isCorrect,
				AnsweredAt:     &now,
			}
			s.repo.CreateAnswer(&answer)
		} else {
			existingAnswer.SelectedOption = selectedOption
			existingAnswer.IsCorrect = isCorrect
			existingAnswer.AnsweredAt = &now
			s.repo.UpdateAnswer(&existingAnswer)
		}
	}

	// Count unanswered
	unansweredCount = len(questions) - correctCount - wrongCount

	// Calculate final score (normalized to MaxScore)
	maxRawScore := float64(subtest.QuestionCount) * WeightHard // Max possible with all hard questions
	finalScore := (rawScore / maxRawScore) * subtest.MaxScore

	// Update subtest result
	now := time.Now()
	result.FinishedAt = &now
	result.CorrectCount = correctCount
	result.WrongCount = wrongCount
	result.UnansweredCount = unansweredCount
	result.RawScore = &rawScore
	result.FinalScore = &finalScore

	if err := s.repo.UpdateSubtestResult(&result); err != nil {
		return nil, err
	}

	utils.LogSuccess("attempts", "submit_subtest", "Subtest submitted", requestID, userID, map[string]any{
		"attempt_id":  attemptID,
		"subtest_id":  subtestID,
		"correct":     correctCount,
		"wrong":       wrongCount,
		"unanswered":  unansweredCount,
		"final_score": finalScore,
	})

	// Fetch updated result
	updatedResult, _ := s.repo.FindSubtestResultByAttemptAndSubtest(attemptID, subtestID)
	response := ToSubtestResultResponse(updatedResult)
	return &response, nil
}

// ==========================================
// Finish Attempt
// ==========================================

func (s *attemptService) FinishAttempt(attemptID uint, userID uint, requestID string) (*AttemptResponse, error) {
	utils.LogInfo("attempts", "finish", "Finishing attempt", requestID, userID, map[string]any{
		"attempt_id": attemptID,
	})

	attempt, err := s.repo.FindAttemptByID(attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("attempt not found")
		}
		return nil, err
	}

	// Check ownership
	if attempt.Registration.UserID != userID {
		return nil, errors.New("you can only finish your own attempt")
	}

	if attempt.Status == entities.AttemptStatusCompleted {
		return nil, errors.New("attempt is already completed")
	}

	// Calculate total score from all subtest results
	results, err := s.repo.FindSubtestResultsByAttemptID(attemptID)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	for _, r := range results {
		if r.FinalScore != nil {
			totalScore += *r.FinalScore
		}
	}
	if len(results) > 0 {
		totalScore = totalScore / float64(len(results))
	}

	now := time.Now()
	attempt.FinishedAt = &now
	attempt.Status = entities.AttemptStatusCompleted
	attempt.TotalScore = &totalScore

	if err := s.repo.UpdateAttempt(&attempt); err != nil {
		return nil, err
	}

	utils.LogSuccess("attempts", "finish", "Attempt finished", requestID, userID, map[string]any{
		"attempt_id":  attemptID,
		"total_score": totalScore,
	})

	// Fetch updated attempt
	finishedAttempt, _ := s.repo.FindAttemptByID(attemptID)
	response := ToAttemptResponse(finishedAttempt)
	return &response, nil
}

// ==========================================
// Leaderboard
// ==========================================

func (s *attemptService) GetLeaderboard(tryOutID uint, requestID string) ([]LeaderboardEntryResponse, error) {
	utils.LogInfo("attempts", "leaderboard", "Fetching leaderboard", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
	})

	attempts, err := s.repo.FindLeaderboard(tryOutID, 100) // Top 100
	if err != nil {
		return nil, err
	}

	var responses []LeaderboardEntryResponse
	for i, a := range attempts {
		score := float64(0)
		if a.TotalScore != nil {
			score = *a.TotalScore
		}

		responses = append(responses, LeaderboardEntryResponse{
			Rank:       i + 1,
			UserID:     a.Registration.UserID,
			Username:   a.Registration.User.Username,
			TotalScore: score,
			FinishedAt: a.FinishedAt,
		})
	}

	return responses, nil
}
