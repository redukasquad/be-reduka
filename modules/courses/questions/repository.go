package questions

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByCourseID(courseID uint) ([]entities.RegistrationQuestion, error)
	FindByID(id uint) (entities.RegistrationQuestion, error)
	Create(question *entities.RegistrationQuestion) error
	Update(question *entities.RegistrationQuestion) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByCourseID(courseID uint) ([]entities.RegistrationQuestion, error) {
	var questions []entities.RegistrationQuestion
	err := r.db.Where("course_id = ?", courseID).Order("question_order ASC").Find(&questions).Error
	return questions, err
}

func (r *repository) FindByID(id uint) (entities.RegistrationQuestion, error) {
	var question entities.RegistrationQuestion
	err := r.db.First(&question, id).Error
	return question, err
}

func (r *repository) Create(question *entities.RegistrationQuestion) error {
	return r.db.Create(question).Error
}

func (r *repository) Update(question *entities.RegistrationQuestion) error {
	return r.db.Save(question).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.RegistrationQuestion{}, id).Error
}
