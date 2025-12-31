package registrations

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindByID(id uint) (entities.CourseRegistration, error)
	FindByUserID(userID uint) ([]entities.CourseRegistration, error)
	FindByCourseID(courseID uint) ([]entities.CourseRegistration, error)
	FindByUserAndCourse(userID, courseID uint) (entities.CourseRegistration, error)
	Create(registration *entities.CourseRegistration) error
	Update(registration *entities.CourseRegistration) error
	Delete(id uint) error
	CreateAnswers(answers []entities.RegistrationAnswer) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(id uint) (entities.CourseRegistration, error) {
	var registration entities.CourseRegistration
	err := r.db.Preload("User").Preload("Course").Preload("Course.Program").Preload("Answers").Preload("Answers.Question").First(&registration, id).Error
	return registration, err
}

func (r *repository) FindByUserID(userID uint) ([]entities.CourseRegistration, error) {
	var registrations []entities.CourseRegistration
	err := r.db.Where("user_id = ?", userID).Preload("User").Preload("Course").Preload("Course.Program").Preload("Answers").Preload("Answers.Question").Find(&registrations).Error
	return registrations, err
}

func (r *repository) FindByCourseID(courseID uint) ([]entities.CourseRegistration, error) {
	var registrations []entities.CourseRegistration
	err := r.db.Where("course_id = ?", courseID).Preload("User").Preload("Course").Preload("Course.Program").Preload("Answers").Preload("Answers.Question").Find(&registrations).Error
	return registrations, err
}

func (r *repository) FindByUserAndCourse(userID, courseID uint) (entities.CourseRegistration, error) {
	var registration entities.CourseRegistration
	err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&registration).Error
	return registration, err
}

func (r *repository) Create(registration *entities.CourseRegistration) error {
	return r.db.Create(registration).Error
}

func (r *repository) Update(registration *entities.CourseRegistration) error {
	return r.db.Save(registration).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.CourseRegistration{}, id).Error
}

func (r *repository) CreateAnswers(answers []entities.RegistrationAnswer) error {
	if len(answers) == 0 {
		return nil
	}
	return r.db.Create(&answers).Error
}
