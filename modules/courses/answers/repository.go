package answers

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByRegistrationID(registrationID uint) ([]entities.RegistrationAnswer, error)
	FindByID(id uint) (entities.RegistrationAnswer, error)
	Create(answer *entities.RegistrationAnswer) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByRegistrationID(registrationID uint) ([]entities.RegistrationAnswer, error) {
	var answers []entities.RegistrationAnswer
	err := r.db.Where("registration_id = ?", registrationID).Preload("Question").Find(&answers).Error
	return answers, err
}

func (r *repository) FindByID(id uint) (entities.RegistrationAnswer, error) {
	var answer entities.RegistrationAnswer
	err := r.db.Preload("Question").First(&answer, id).Error
	return answer, err
}

func (r *repository) Create(answer *entities.RegistrationAnswer) error {
	return r.db.Create(answer).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.RegistrationAnswer{}, id).Error
}
