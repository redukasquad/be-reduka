package questions

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// Subtests
	FindAllSubtests() ([]entities.Subtest, error)
	FindSubtestByID(id uint) (entities.Subtest, error)

	// Questions
	FindByTryOutID(tryOutID uint) ([]entities.TryOutQuestion, error)
	FindByTryOutAndSubtest(tryOutID, subtestID uint) ([]entities.TryOutQuestion, error)
	FindByID(id uint) (entities.TryOutQuestion, error)
	CountByTryOutAndSubtest(tryOutID, subtestID uint) (int64, error)
	Create(question *entities.TryOutQuestion) error
	Update(question *entities.TryOutQuestion) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==========================================
// Subtest Repository Methods
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
// Question Repository Methods
// ==========================================

func (r *repository) FindByTryOutID(tryOutID uint) ([]entities.TryOutQuestion, error) {
	var questions []entities.TryOutQuestion
	err := r.db.Where("try_out_package_id = ?", tryOutID).
		Preload("Subtest").
		Order("subtest_id ASC, order_number ASC").
		Find(&questions).Error
	return questions, err
}

func (r *repository) FindByTryOutAndSubtest(tryOutID, subtestID uint) ([]entities.TryOutQuestion, error) {
	var questions []entities.TryOutQuestion
	err := r.db.Where("try_out_package_id = ? AND subtest_id = ?", tryOutID, subtestID).
		Preload("Subtest").
		Order("order_number ASC").
		Find(&questions).Error
	return questions, err
}

func (r *repository) FindByID(id uint) (entities.TryOutQuestion, error) {
	var question entities.TryOutQuestion
	err := r.db.Preload("Subtest").First(&question, id).Error
	return question, err
}

func (r *repository) CountByTryOutAndSubtest(tryOutID, subtestID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.TryOutQuestion{}).
		Where("try_out_package_id = ? AND subtest_id = ?", tryOutID, subtestID).
		Count(&count).Error
	return count, err
}

func (r *repository) Create(question *entities.TryOutQuestion) error {
	return r.db.Create(question).Error
}

func (r *repository) Update(question *entities.TryOutQuestion) error {
	return r.db.Save(question).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.TryOutQuestion{}, id).Error
}
