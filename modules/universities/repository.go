package universities

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	// University
	FindAllUniversities(search string) ([]entities.University, error)
	FindUniversityByID(id uint) (entities.University, error)
	CreateUniversity(u *entities.University) error
	UpdateUniversity(u entities.University) error
	DeleteUniversity(id uint) error

	// Major
	FindMajorsByUniversity(universityID uint) ([]entities.UniversityMajor, error)
	FindMajorByID(id uint) (entities.UniversityMajor, error)
	CreateMajor(m *entities.UniversityMajor) error
	UpdateMajor(m entities.UniversityMajor) error
	DeleteMajor(id uint) error

	// UserTarget
	FindTargetsByUser(userID uint) ([]entities.UserTarget, error)
	FindTargetByUserAndMajor(userID, majorID uint) (entities.UserTarget, error)
	FindUsersByUniversity(universityID uint) ([]entities.User, error)
	CreateTarget(t *entities.UserTarget) error
	UpdateTarget(t entities.UserTarget) error
	DeleteTarget(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllUniversities(search string) ([]entities.University, error) {
	var unis []entities.University
	q := r.db.Preload("Major")
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	err := q.Find(&unis).Error
	return unis, err
}

func (r *repository) FindUniversityByID(id uint) (entities.University, error) {
	var u entities.University
	err := r.db.Preload("Major").First(&u, id).Error
	return u, err
}

func (r *repository) CreateUniversity(u *entities.University) error {
	return r.db.Create(u).Error
}

func (r *repository) UpdateUniversity(u entities.University) error {
	return r.db.Save(&u).Error
}

func (r *repository) DeleteUniversity(id uint) error {
	return r.db.Delete(&entities.University{}, id).Error
}

func (r *repository) FindMajorsByUniversity(universityID uint) ([]entities.UniversityMajor, error) {
	var majors []entities.UniversityMajor
	err := r.db.Where("university_id = ?", universityID).Preload("University").Find(&majors).Error
	return majors, err
}

func (r *repository) FindMajorByID(id uint) (entities.UniversityMajor, error) {
	var m entities.UniversityMajor
	err := r.db.Preload("University").First(&m, id).Error
	return m, err
}

func (r *repository) CreateMajor(m *entities.UniversityMajor) error {
	return r.db.Create(m).Error
}

func (r *repository) UpdateMajor(m entities.UniversityMajor) error {
	return r.db.Save(&m).Error
}

func (r *repository) DeleteMajor(id uint) error {
	return r.db.Delete(&entities.UniversityMajor{}, id).Error
}

func (r *repository) FindTargetsByUser(userID uint) ([]entities.UserTarget, error) {
	var targets []entities.UserTarget
	err := r.db.Where("user_id = ?", userID).
		Preload("Major").
		Preload("Major.University").
		Order("priority asc").
		Find(&targets).Error
	return targets, err
}

func (r *repository) FindTargetByUserAndMajor(userID, majorID uint) (entities.UserTarget, error) {
	var t entities.UserTarget
	err := r.db.Where("user_id = ? AND university_major_id = ?", userID, majorID).First(&t).Error
	return t, err
}

func (r *repository) FindUsersByUniversity(universityID uint) ([]entities.User, error) {
	var users []entities.User
	err := r.db.
		Joins("JOIN user_targets ON user_targets.user_id = users.id").
		Joins("JOIN university_majors ON university_majors.id = user_targets.university_major_id").
		Where("university_majors.university_id = ? AND user_targets.deleted_at IS NULL AND university_majors.deleted_at IS NULL", universityID).
		Distinct("users.id").
		Find(&users).Error
	return users, err
}

func (r *repository) CreateTarget(t *entities.UserTarget) error {
	return r.db.Create(t).Error
}

func (r *repository) UpdateTarget(t entities.UserTarget) error {
	return r.db.Save(&t).Error
}

func (r *repository) DeleteTarget(id uint) error {
	return r.db.Delete(&entities.UserTarget{}, id).Error
}
