package subjects

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByCourseID(courseID uint) ([]entities.Class, error)
	FindByID(id uint) (entities.Class, error)
	Create(class *entities.Class) error
	Update(class *entities.Class) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByCourseID(courseID uint) ([]entities.Class, error) {
	var classes []entities.Class
	err := r.db.Where("course_id = ?", courseID).Preload("Lessons").Find(&classes).Error
	return classes, err
}

func (r *repository) FindByID(id uint) (entities.Class, error) {
	var class entities.Class
	err := r.db.Preload("Course").Preload("Lessons").First(&class, id).Error
	return class, err
}

func (r *repository) Create(class *entities.Class) error {
	return r.db.Create(class).Error
}

func (r *repository) Update(class *entities.Class) error {
	return r.db.Save(class).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.Class{}, id).Error
}
