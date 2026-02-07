package attempts

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// Registration
	FindRegistrationByID(id uint) (entities.TryOutRegistration, error)
	FindRegistrationByUserAndTryOut(userID, tryOutID uint) (entities.TryOutRegistration, error)

	// Attempt
	FindAttemptByID(id uint) (entities.TryOutAttempt, error)
	FindAttemptByRegistrationID(registrationID uint) (entities.TryOutAttempt, error)
	CreateAttempt(attempt *entities.TryOutAttempt) error
	UpdateAttempt(attempt *entities.TryOutAttempt) error

	// Subtest Result
	FindSubtestResultByAttemptAndSubtest(attemptID, subtestID uint) (entities.SubtestResult, error)
	FindSubtestResultsByAttemptID(attemptID uint) ([]entities.SubtestResult, error)
	CreateSubtestResult(result *entities.SubtestResult) error
	UpdateSubtestResult(result *entities.SubtestResult) error

	// Answers
	FindAnswerByAttemptAndQuestion(attemptID, questionID uint) (entities.UserTryOutAnswer, error)
	FindAnswersByAttemptAndSubtest(attemptID, subtestID uint) ([]entities.UserTryOutAnswer, error)
	CreateAnswer(answer *entities.UserTryOutAnswer) error
	UpdateAnswer(answer *entities.UserTryOutAnswer) error

	// Questions
	FindQuestionsByTryOutAndSubtest(tryOutID, subtestID uint) ([]entities.TryOutQuestion, error)
	FindQuestionByID(id uint) (entities.TryOutQuestion, error)

	// Subtests
	FindAllSubtests() ([]entities.Subtest, error)
	FindSubtestByID(id uint) (entities.Subtest, error)

	// Leaderboard
	FindLeaderboard(tryOutID uint, limit int) ([]entities.TryOutAttempt, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==========================================
// Registration Methods
// ==========================================

func (r *repository) FindRegistrationByID(id uint) (entities.TryOutRegistration, error) {
	var registration entities.TryOutRegistration
	err := r.db.Preload("TryOutPackage").Preload("User").First(&registration, id).Error
	return registration, err
}

func (r *repository) FindRegistrationByUserAndTryOut(userID, tryOutID uint) (entities.TryOutRegistration, error) {
	var registration entities.TryOutRegistration
	err := r.db.Where("user_id = ? AND try_out_package_id = ?", userID, tryOutID).
		Preload("TryOutPackage").
		First(&registration).Error
	return registration, err
}

// ==========================================
// Attempt Methods
// ==========================================

func (r *repository) FindAttemptByID(id uint) (entities.TryOutAttempt, error) {
	var attempt entities.TryOutAttempt
	err := r.db.Preload("Registration.TryOutPackage").
		Preload("Registration.User").
		Preload("CurrentSubtest").
		Preload("SubtestResults.Subtest").
		First(&attempt, id).Error
	return attempt, err
}

func (r *repository) FindAttemptByRegistrationID(registrationID uint) (entities.TryOutAttempt, error) {
	var attempt entities.TryOutAttempt
	err := r.db.Where("registration_id = ?", registrationID).
		Preload("Registration.TryOutPackage").
		Preload("Registration.User").
		Preload("CurrentSubtest").
		Preload("SubtestResults.Subtest").
		First(&attempt).Error
	return attempt, err
}

func (r *repository) CreateAttempt(attempt *entities.TryOutAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *repository) UpdateAttempt(attempt *entities.TryOutAttempt) error {
	return r.db.Save(attempt).Error
}

// ==========================================
// Subtest Result Methods
// ==========================================

func (r *repository) FindSubtestResultByAttemptAndSubtest(attemptID, subtestID uint) (entities.SubtestResult, error) {
	var result entities.SubtestResult
	err := r.db.Where("attempt_id = ? AND subtest_id = ?", attemptID, subtestID).
		Preload("Subtest").
		First(&result).Error
	return result, err
}

func (r *repository) FindSubtestResultsByAttemptID(attemptID uint) ([]entities.SubtestResult, error) {
	var results []entities.SubtestResult
	err := r.db.Where("attempt_id = ?", attemptID).
		Preload("Subtest").
		Order("subtest_id ASC").
		Find(&results).Error
	return results, err
}

func (r *repository) CreateSubtestResult(result *entities.SubtestResult) error {
	return r.db.Create(result).Error
}

func (r *repository) UpdateSubtestResult(result *entities.SubtestResult) error {
	return r.db.Save(result).Error
}

// ==========================================
// Answer Methods
// ==========================================

func (r *repository) FindAnswerByAttemptAndQuestion(attemptID, questionID uint) (entities.UserTryOutAnswer, error) {
	var answer entities.UserTryOutAnswer
	err := r.db.Where("attempt_id = ? AND question_id = ?", attemptID, questionID).First(&answer).Error
	return answer, err
}

func (r *repository) FindAnswersByAttemptAndSubtest(attemptID, subtestID uint) ([]entities.UserTryOutAnswer, error) {
	var answers []entities.UserTryOutAnswer
	err := r.db.Joins("JOIN try_out_questions ON try_out_questions.id = user_try_out_answers.question_id").
		Where("user_try_out_answers.attempt_id = ? AND try_out_questions.subtest_id = ?", attemptID, subtestID).
		Find(&answers).Error
	return answers, err
}

func (r *repository) CreateAnswer(answer *entities.UserTryOutAnswer) error {
	return r.db.Create(answer).Error
}

func (r *repository) UpdateAnswer(answer *entities.UserTryOutAnswer) error {
	return r.db.Save(answer).Error
}

// ==========================================
// Question Methods
// ==========================================

func (r *repository) FindQuestionsByTryOutAndSubtest(tryOutID, subtestID uint) ([]entities.TryOutQuestion, error) {
	var questions []entities.TryOutQuestion
	err := r.db.Where("try_out_package_id = ? AND subtest_id = ?", tryOutID, subtestID).
		Order("order_number ASC").
		Find(&questions).Error
	return questions, err
}

func (r *repository) FindQuestionByID(id uint) (entities.TryOutQuestion, error) {
	var question entities.TryOutQuestion
	err := r.db.First(&question, id).Error
	return question, err
}

// ==========================================
// Subtest Methods
// ==========================================

func (r *repository) FindAllSubtests() ([]entities.Subtest, error) {
	var subtests []entities.Subtest
	err := r.db.Order("id ASC").Find(&subtests).Error
	return subtests, err
}

func (r *repository) FindSubtestByID(id uint) (entities.Subtest, error) {
	var subtest entities.Subtest
	err := r.db.First(&subtest, id).Error
	return subtest, err
}

// ==========================================
// Leaderboard Methods
// ==========================================

func (r *repository) FindLeaderboard(tryOutID uint, limit int) ([]entities.TryOutAttempt, error) {
	var attempts []entities.TryOutAttempt
	err := r.db.Joins("JOIN try_out_registrations ON try_out_registrations.id = try_out_attempts.registration_id").
		Where("try_out_registrations.try_out_package_id = ?", tryOutID).
		Where("try_out_attempts.status = ?", "completed").
		Where("try_out_attempts.total_score IS NOT NULL").
		Preload("Registration.User").
		Order("try_out_attempts.total_score DESC, try_out_attempts.finished_at ASC").
		Limit(limit).
		Find(&attempts).Error
	return attempts, err
}
