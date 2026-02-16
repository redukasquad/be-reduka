package uploads

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(image *entities.Image) error
	FindByURL(url string) (entities.Image, error)
	DeleteByURL(url string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(image *entities.Image) error {
	return r.db.Create(image).Error
}

func (r *repository) FindByURL(url string) (entities.Image, error) {
	var image entities.Image
	err := r.db.Where("url = ?", url).First(&image).Error
	return image, err
}

func (r *repository) DeleteByURL(url string) error {
	return r.db.Where("url = ?", url).Delete(&entities.Image{}).Error
}
