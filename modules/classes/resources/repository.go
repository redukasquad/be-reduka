package resources

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByLessonID(lessonID uint) ([]entities.LessonResource, error)
	FindByID(id uint) (entities.LessonResource, error)
	Create(resource *entities.LessonResource) error
	Update(resource *entities.LessonResource) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByLessonID(lessonID uint) ([]entities.LessonResource, error) {
	var resources []entities.LessonResource
	err := r.db.Where("lesson_id = ?", lessonID).Find(&resources).Error
	return resources, err
}

func (r *repository) FindByID(id uint) (entities.LessonResource, error) {
	var resource entities.LessonResource
	err := r.db.Preload("Lesson").First(&resource, id).Error
	return resource, err
}

func (r *repository) Create(resource *entities.LessonResource) error {
	return r.db.Create(resource).Error
}

func (r *repository) Update(resource *entities.LessonResource) error {
	return r.db.Save(resource).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.LessonResource{}, id).Error
}
