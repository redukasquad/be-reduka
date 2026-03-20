package courses

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll() ([]entities.Course, error)
	FindAllPaginated(offset, limit int, search string) ([]entities.Course, error)
	FindByCreatorPaginated(creatorID uint, offset, limit int, search string) ([]entities.Course, error)
	CountByCreator(creatorID uint, search string) (int64, error)
	CountWithSearch(search string) (int64, error)
	FindByID(id uint) (entities.Course, error)
	FindByProgramID(programID uint) ([]entities.Course, error)
	FindByName(name string) (entities.Course, error)
	Create(course *entities.Course) error
	Update(course *entities.Course) error
	Delete(id uint) error
	Count() (int64, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByCreatorPaginated(creatorID uint, offset, limit int, search string) ([]entities.Course, error) {
	var courses []entities.Course
	query := r.db.Preload("Program").Preload("Creator").Preload("Classes").
		Where("created_by_user_id = ?", creatorID)
	if search != "" {
		query = query.Where("name_course LIKE ?", "%"+search+"%")
	}
	err := query.Offset(offset).Limit(limit).Find(&courses).Error
	return courses, err
}

func (r *repository) CountByCreator(creatorID uint, search string) (int64, error) {
	var count int64
	query := r.db.Model(&entities.Course{}).Where("created_by_user_id = ?", creatorID)
	if search != "" {
		query = query.Where("name_course LIKE ?", "%"+search+"%")
	}
	return count, query.Count(&count).Error
}

func (r *repository) FindAll() ([]entities.Course, error) {
	var courses []entities.Course
	err := r.db.Preload("Program").Preload("Creator").Find(&courses).Error
	return courses, err
}

func (r *repository) FindAllPaginated(offset, limit int, search string) ([]entities.Course, error) {
	var course []entities.Course

	query := r.db.Preload("Program").Preload("Creator").Preload("Classes")

	if search != "" {
		query = query.Where("name_course LIKE ?", "%"+search+"%")
	}

	err := query.Offset(offset).Limit(limit).Find(&course).Error
	return course, err
}

func (r *repository) CountWithSearch(search string) (int64, error) {
	var count int64

	query := r.db.Model(&entities.Course{})

	if search != "" {
		query = query.Where("name_course LIKE ?", "%"+search+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *repository) FindByID(id uint) (entities.Course, error) {
	var course entities.Course
	err := r.db.Preload("Program").Preload("Creator").Preload("Classes").First(&course, id).Error
	return course, err
}

func (r *repository) FindByProgramID(programID uint) ([]entities.Course, error) {
	var courses []entities.Course
	err := r.db.Where("program_id = ?", programID).Preload("Program").Preload("Creator").Find(&courses).Error
	return courses, err
}

func (r *repository) FindByName(name string) (entities.Course, error) {
	var course entities.Course
	err := r.db.Where("name_course = ?", name).Preload("Program").First(&course).Error
	return course, err
}

func (r *repository) Create(course *entities.Course) error {
	return r.db.Create(course).Error
}

func (r *repository) Update(course *entities.Course) error {
	return r.db.Save(course).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.Course{}, id).Error
}

func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.Course{}).Count(&count).Error
	return count, err
}
