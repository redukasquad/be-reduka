package courses

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Repository interface defines the methods for course data access
type Repository interface {
	FindAll() ([]entities.Course, error)
	FindByID(id uint) (entities.Course, error)
	FindByProgramID(programID uint) ([]entities.Course, error)
	FindByName(name string) (entities.Course, error)
	Create(course *entities.Course) error
	Update(course *entities.Course) error
	Delete(id uint) error
	Count() (int64, error)
}

// NewRepository creates a new course repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]entities.Course, error) {
	var courses []entities.Course
	err := r.db.Preload("Program").Preload("Creator").Find(&courses).Error
	return courses, err
}

func (r *repository) FindByID(id uint) (entities.Course, error) {
	var course entities.Course
	err := r.db.Preload("Program").Preload("Creator").Preload("Subjects").First(&course, id).Error
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
