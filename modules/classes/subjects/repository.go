package subjects

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByCourseID(courseID uint) ([]entities.ClassSubject, error)
	FindByID(id uint) (entities.ClassSubject, error)
	Create(subject *entities.ClassSubject) error
	Update(subject *entities.ClassSubject) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByCourseID(courseID uint) ([]entities.ClassSubject, error) {
	var subjects []entities.ClassSubject
	err := r.db.Where("course_id = ?", courseID).Preload("Lessons").Find(&subjects).Error
	return subjects, err
}

func (r *repository) FindByID(id uint) (entities.ClassSubject, error) {
	var subject entities.ClassSubject
	err := r.db.Preload("Course").Preload("Lessons").First(&subject, id).Error
	return subject, err
}

func (r *repository) Create(subject *entities.ClassSubject) error {
	return r.db.Create(subject).Error
}

func (r *repository) Update(subject *entities.ClassSubject) error {
	return r.db.Save(subject).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.ClassSubject{}, id).Error
}
