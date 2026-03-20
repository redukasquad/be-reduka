package lessons

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByClassID(classID uint) ([]entities.Lesson, error)
	FindByID(id uint) (entities.Lesson, error)
	Create(lesson *entities.Lesson) error
	Update(lesson *entities.Lesson) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByClassID(classID uint) ([]entities.Lesson, error) {
	var lessons []entities.Lesson
	err := r.db.Where("class_id = ?", classID).Order("lesson_order ASC").Preload("Resources").Find(&lessons).Error
	return lessons, err
}

func (r *repository) FindByID(id uint) (entities.Lesson, error) {
	var lesson entities.Lesson
	err := r.db.Preload("Class").Preload("Resources").Preload("Creator").First(&lesson, id).Error
	return lesson, err
}

func (r *repository) Create(lesson *entities.Lesson) error {
	return r.db.Create(lesson).Error
}

func (r *repository) Update(lesson *entities.Lesson) error {
	return r.db.Save(lesson).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.Lesson{}, id).Error
}
