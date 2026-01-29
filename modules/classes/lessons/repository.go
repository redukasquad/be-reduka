package lessons

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindBySubjectID(subjectID uint) ([]entities.ClassLesson, error)
	FindByID(id uint) (entities.ClassLesson, error)
	Create(lesson *entities.ClassLesson) error
	Update(lesson *entities.ClassLesson) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindBySubjectID(subjectID uint) ([]entities.ClassLesson, error) {
	var lessons []entities.ClassLesson
	err := r.db.Where("subject_id = ?", subjectID).Order("lesson_order ASC").Preload("Resources").Find(&lessons).Error
	return lessons, err
}

func (r *repository) FindByID(id uint) (entities.ClassLesson, error) {
	var lesson entities.ClassLesson
	err := r.db.Preload("Subject").Preload("Resources").Preload("Creator").First(&lesson, id).Error
	return lesson, err
}

func (r *repository) Create(lesson *entities.ClassLesson) error {
	return r.db.Create(lesson).Error
}

func (r *repository) Update(lesson *entities.ClassLesson) error {
	return r.db.Save(lesson).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.ClassLesson{}, id).Error
}
