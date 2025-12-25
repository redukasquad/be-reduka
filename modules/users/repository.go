package users

import (
	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	FindAll() ([]entities.User, error)
	FindByID(id int) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	FindByVerificationCode(code string) (entities.User, error)
	FindByResetToken(token string) (entities.User, error)
	Create(user *entities.User) error
	Update(user entities.User) error
	Delete(id int) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(id int) (entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *repository) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *repository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByVerificationCode(code string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("verification_code = ?", code).First(&user).Error
	return user, err
}

func (r *repository) Update(user entities.User) error {
	return r.db.Save(&user).Error
}

func (r *repository) FindByResetToken(token string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("reset_password_token = ?", token).First(&user).Error
	return user, err
}

func (r *repository) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) Delete(id int) error {
	return r.db.Delete(&entities.User{}, id).Error
}
