package users

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type userService struct {
	repo Repository
}

type Service interface {
	GetAll() ([]entities.User, error)
	GetByID(id int) (*entities.User, error)
	Update(id int, input UpdateUserInput) (*entities.User, error)
	SetRole(userID int, role string) (*entities.User, error)
	Delete(id int) error
}

func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]entities.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id int) (*entities.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.Password = "" // Don't expose password
	return &user, nil
}

func (s *userService) Update(id int, input UpdateUserInput) (*entities.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update only provided fields (Role is NOT allowed here - use SetRole instead)
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.NoTelp != nil {
		user.NoTelp = *input.NoTelp
	}
	if input.JenisKelamin != nil {
		user.JenisKelamin = input.JenisKelamin
	}
	if input.Kelas != nil {
		user.Kelas = input.Kelas
	}
	if input.ProfileImage != nil {
		user.ProfileImage = *input.ProfileImage
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	user.Password = "" // Don't expose password
	return &user, nil
}

// SetRole allows admin to set a user's role
func (s *userService) SetRole(userID int, role string) (*entities.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Role = &role

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	user.Password = "" // Don't expose password
	return &user, nil
}

func (s *userService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
