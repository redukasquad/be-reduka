package registrations

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// Registrations
	FindByID(id uint) (entities.TryOutRegistration, error)
	FindByUserAndTryOut(userID, tryOutID uint) (entities.TryOutRegistration, error)
	FindByUserID(userID uint) ([]entities.TryOutRegistration, error)
	FindPendingPayments() ([]entities.TryOutRegistration, error)
	FindByTryOutID(tryOutID uint) ([]entities.TryOutRegistration, error)
	Create(registration *entities.TryOutRegistration) error
	Update(registration *entities.TryOutRegistration) error

	// Try Out
	FindTryOutByID(id uint) (entities.TryOut, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==========================================
// Registration Repository Methods
// ==========================================

func (r *repository) FindByID(id uint) (entities.TryOutRegistration, error) {
	var registration entities.TryOutRegistration
	err := r.db.Preload("TryOut").
		Preload("User").
		Preload("ApprovedBy").
		Preload("Attempt").
		First(&registration, id).Error
	return registration, err
}

func (r *repository) FindByUserAndTryOut(userID, tryOutID uint) (entities.TryOutRegistration, error) {
	var registration entities.TryOutRegistration
	err := r.db.Where("user_id = ? AND try_out_package_id = ?", userID, tryOutID).
		Preload("TryOut").
		Preload("User").
		Preload("ApprovedBy").
		Preload("Attempt").
		First(&registration).Error
	return registration, err
}

func (r *repository) FindByUserID(userID uint) ([]entities.TryOutRegistration, error) {
	var registrations []entities.TryOutRegistration
	err := r.db.Where("user_id = ?", userID).
		Preload("TryOut").
		Preload("User").
		Preload("ApprovedBy").
		Preload("Attempt").
		Order("registered_at DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *repository) FindPendingPayments() ([]entities.TryOutRegistration, error) {
	var registrations []entities.TryOutRegistration
	err := r.db.Where("payment_status = ?", "pending").
		Where("payment_proof_url IS NOT NULL AND payment_proof_url != ''").
		Preload("TryOut").
		Preload("User").
		Order("registered_at ASC").
		Find(&registrations).Error
	return registrations, err
}

func (r *repository) FindByTryOutID(tryOutID uint) ([]entities.TryOutRegistration, error) {
	var registrations []entities.TryOutRegistration
	err := r.db.Where("try_out_package_id = ?", tryOutID).
		Preload("TryOut").
		Preload("User").
		Preload("ApprovedBy").
		Order("registered_at DESC").
		Find(&registrations).Error
	return registrations, err
}

func (r *repository) Create(registration *entities.TryOutRegistration) error {
	return r.db.Create(registration).Error
}

func (r *repository) Update(registration *entities.TryOutRegistration) error {
	return r.db.Save(registration).Error
}

// ==========================================
// Try Out Repository Methods
// ==========================================

func (r *repository) FindTryOutByID(id uint) (entities.TryOut, error) {
	var tryOut entities.TryOut
	err := r.db.First(&tryOut, id).Error
	return tryOut, err
}
