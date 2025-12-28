package answers

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Repository interface defines the methods for answer data access
type Repository interface {
	FindByRegistrationID(registrationID uint) ([]entities.RegistrationAnswer, error)
}

// NewRepository creates a new answer repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByRegistrationID(registrationID uint) ([]entities.RegistrationAnswer, error) {
	var answers []entities.RegistrationAnswer
	err := r.db.Where("registration_id = ?", registrationID).Preload("Question").Find(&answers).Error
	return answers, err
}
