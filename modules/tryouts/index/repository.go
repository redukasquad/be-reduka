package tryouts

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// Try Out
	FindAll() ([]entities.TryOut, error)
	FindAllPaginated(offset, limit int, search string, publishedOnly bool) ([]entities.TryOut, error)
	CountWithSearch(search string, publishedOnly bool) (int64, error)
	FindByID(id uint) (entities.TryOut, error)
	FindByName(name string) (entities.TryOut, error)
	Create(tryOut *entities.TryOut) error
	Update(tryOut *entities.TryOut) error
	Delete(id uint) error

	// Tutor Permissions
	FindTutorPermissions(tryOutID uint) ([]entities.TutorPermission, error)
	FindTutorPermission(tryOutID, userID uint) (entities.TutorPermission, error)
	CreateTutorPermission(permission *entities.TutorPermission) error
	DeleteTutorPermission(tryOutID, userID uint) error
	HasTutorPermission(tryOutID, userID uint) (bool, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ==========================================
// Try Out Repository Methods
// ==========================================

func (r *repository) FindAll() ([]entities.TryOut, error) {
	var tryOuts []entities.TryOut
	err := r.db.Preload("Creator").Find(&tryOuts).Error
	return tryOuts, err
}

func (r *repository) FindAllPaginated(offset, limit int, search string, publishedOnly bool) ([]entities.TryOut, error) {
	var tryOuts []entities.TryOut

	query := r.db.Preload("Creator")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tryOuts).Error
	return tryOuts, err
}

func (r *repository) CountWithSearch(search string, publishedOnly bool) (int64, error) {
	var count int64

	query := r.db.Model(&entities.TryOut{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *repository) FindByID(id uint) (entities.TryOut, error) {
	var tryOut entities.TryOut
	err := r.db.Preload("Creator").First(&tryOut, id).Error
	return tryOut, err
}

func (r *repository) FindByName(name string) (entities.TryOut, error) {
	var tryOut entities.TryOut
	err := r.db.Where("name = ?", name).First(&tryOut).Error
	return tryOut, err
}

func (r *repository) Create(tryOut *entities.TryOut) error {
	return r.db.Create(tryOut).Error
}

func (r *repository) Update(tryOut *entities.TryOut) error {
	return r.db.Save(tryOut).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.TryOut{}, id).Error
}

// ==========================================
// Tutor Permission Repository Methods
// ==========================================

func (r *repository) FindTutorPermissions(tryOutID uint) ([]entities.TutorPermission, error) {
	var permissions []entities.TutorPermission
	err := r.db.Where("try_out_package_id = ?", tryOutID).
		Preload("User").
		Preload("GrantedBy").
		Find(&permissions).Error
	return permissions, err
}

func (r *repository) FindTutorPermission(tryOutID, userID uint) (entities.TutorPermission, error) {
	var permission entities.TutorPermission
	err := r.db.Where("try_out_package_id = ? AND user_id = ?", tryOutID, userID).
		Preload("User").
		Preload("GrantedBy").
		First(&permission).Error
	return permission, err
}

func (r *repository) CreateTutorPermission(permission *entities.TutorPermission) error {
	return r.db.Create(permission).Error
}

func (r *repository) DeleteTutorPermission(tryOutID, userID uint) error {
	return r.db.Where("try_out_package_id = ? AND user_id = ?", tryOutID, userID).
		Delete(&entities.TutorPermission{}).Error
}

func (r *repository) HasTutorPermission(tryOutID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entities.TutorPermission{}).
		Where("try_out_package_id = ? AND user_id = ?", tryOutID, userID).
		Count(&count).Error
	return count > 0, err
}
