package resources

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByLessonID(lessonID uint) ([]entities.ClassLessonResource, error)
	FindByID(id uint) (entities.ClassLessonResource, error)
	Create(resource *entities.ClassLessonResource) error
	Update(resource *entities.ClassLessonResource) error
	Delete(id uint) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByLessonID(lessonID uint) ([]entities.ClassLessonResource, error) {
	var resources []entities.ClassLessonResource
	err := r.db.Where("class_lesson_id = ?", lessonID).Find(&resources).Error
	return resources, err
}

func (r *repository) FindByID(id uint) (entities.ClassLessonResource, error) {
	var resource entities.ClassLessonResource
	err := r.db.Preload("ClassLesson").First(&resource, id).Error
	return resource, err
}

func (r *repository) Create(resource *entities.ClassLessonResource) error {
	return r.db.Create(resource).Error
}

func (r *repository) Update(resource *entities.ClassLessonResource) error {
	return r.db.Save(resource).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.ClassLessonResource{}, id).Error
}
