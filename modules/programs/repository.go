package programs

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// Repository interface defines the methods for program data access
type Repository interface {
	FindAll() ([]entities.Program, error)
	FindByID(id uint) (entities.Program, error)
	FindByName(name string) (entities.Program, error)
	Create(program *entities.Program) error
	Update(program entities.Program) error
	Delete(id uint) error
	Count() (int64, error)
}

// NewRepository creates a new program repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]entities.Program, error) {
	var programs []entities.Program
	err := r.db.Preload("Courses").Find(&programs).Error
	return programs, err
}

func (r *repository) FindByID(id uint) (entities.Program, error) {
	var program entities.Program
	err := r.db.Preload("Courses").First(&program, id).Error
	return program, err
}

func (r *repository) FindByName(name string) (entities.Program, error) {
	var program entities.Program
	err := r.db.Where("program_name = ?", name).Preload("Courses").First(&program).Error
	return program, err
}

func (r *repository) Create(program *entities.Program) error {
	return r.db.Create(program).Error
}

func (r *repository) Update(program entities.Program) error {
	return r.db.Save(&program).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&entities.Program{}, id).Error
}

func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.Program{}).Count(&count).Error
	return count, err
}
