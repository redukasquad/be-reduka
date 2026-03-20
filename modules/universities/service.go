package universities

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

type universityService struct {
	repo Repository
}

type Service interface {
	GetAllUniversities(search string) ([]UniversityResponse, error)
	GetUniversityByID(id uint) (*UniversityResponse, error)
	CreateUniversity(input CreateUniversityInput) (*UniversityResponse, error)
	UpdateUniversity(id uint, input UpdateUniversityInput) (*UniversityResponse, error)
	DeleteUniversity(id uint) error

	GetMajorsByUniversity(universityID uint) ([]MajorResponse, error)
	CreateMajor(input CreateMajorInput) (*MajorResponse, error)
	UpdateMajor(id uint, input UpdateMajorInput) (*MajorResponse, error)
	DeleteMajor(id uint) error

	GetMyTargets(userID uint) ([]UserTargetResponse, error)
	AddTarget(userID uint, input SetUserTargetInput) (*UserTargetResponse, error)
	DeleteTarget(targetID uint, userID uint) error
	GetUsersByUniversity(universityID uint) ([]UserResponse, error)
}

func NewService(repo Repository) Service {
	return &universityService{repo: repo}
}

func (s *universityService) GetAllUniversities(search string) ([]UniversityResponse, error) {
	unis, err := s.repo.FindAllUniversities(search)
	if err != nil {
		return nil, err
	}
	result := make([]UniversityResponse, len(unis))
	for i, u := range unis {
		result[i] = ToUniversityResponse(u)
	}
	return result, nil
}

func (s *universityService) GetUniversityByID(id uint) (*UniversityResponse, error) {
	u, err := s.repo.FindUniversityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("university not found")
		}
		return nil, err
	}
	resp := ToUniversityResponse(u)
	return &resp, nil
}

func (s *universityService) CreateUniversity(input CreateUniversityInput) (*UniversityResponse, error) {
	u := &entities.University{Name: input.Name, Type: input.Type}
	if err := s.repo.CreateUniversity(u); err != nil {
		return nil, err
	}
	resp := ToUniversityResponse(*u)
	return &resp, nil
}

func (s *universityService) UpdateUniversity(id uint, input UpdateUniversityInput) (*UniversityResponse, error) {
	u, err := s.repo.FindUniversityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("university not found")
		}
		return nil, err
	}
	if input.Name != nil {
		u.Name = *input.Name
	}
	if input.Type != nil {
		u.Type = *input.Type
	}
	if err := s.repo.UpdateUniversity(u); err != nil {
		return nil, err
	}
	resp := ToUniversityResponse(u)
	return &resp, nil
}

func (s *universityService) DeleteUniversity(id uint) error {
	_, err := s.repo.FindUniversityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("university not found")
		}
		return err
	}
	return s.repo.DeleteUniversity(id)
}

func (s *universityService) GetMajorsByUniversity(universityID uint) ([]MajorResponse, error) {
	majors, err := s.repo.FindMajorsByUniversity(universityID)
	if err != nil {
		return nil, err
	}
	result := make([]MajorResponse, len(majors))
	for i, m := range majors {
		result[i] = ToMajorResponse(m)
	}
	return result, nil
}

func (s *universityService) CreateMajor(input CreateMajorInput) (*MajorResponse, error) {
	m := &entities.UniversityMajor{
		UniversityID: input.UniversityID,
		Name:         input.Name,
		PassingGrade: input.PassingGrade,
	}
	if err := s.repo.CreateMajor(m); err != nil {
		return nil, err
	}
	resp := ToMajorResponse(*m)
	return &resp, nil
}

func (s *universityService) UpdateMajor(id uint, input UpdateMajorInput) (*MajorResponse, error) {
	m, err := s.repo.FindMajorByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("major not found")
		}
		return nil, err
	}
	if input.Name != nil {
		m.Name = *input.Name
	}
	if input.PassingGrade != nil {
		m.PassingGrade = *input.PassingGrade
	}
	if err := s.repo.UpdateMajor(m); err != nil {
		return nil, err
	}
	resp := ToMajorResponse(m)
	return &resp, nil
}

func (s *universityService) DeleteMajor(id uint) error {
	_, err := s.repo.FindMajorByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("major not found")
		}
		return err
	}
	return s.repo.DeleteMajor(id)
}

func (s *universityService) GetMyTargets(userID uint) ([]UserTargetResponse, error) {
	targets, err := s.repo.FindTargetsByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]UserTargetResponse, len(targets))
	for i, t := range targets {
		result[i] = ToUserTargetResponse(t)
	}
	return result, nil
}

func (s *universityService) AddTarget(userID uint, input SetUserTargetInput) (*UserTargetResponse, error) {
	// Check major exists
	_, err := s.repo.FindMajorByID(input.UniversityMajorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("major not found")
		}
		return nil, err
	}

	// Check not already added
	existing, err := s.repo.FindTargetByUserAndMajor(userID, input.UniversityMajorID)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("major already in your target list")
	}

	t := &entities.UserTarget{
		UserID:            userID,
		UniversityMajorID: input.UniversityMajorID,
		Priority:          input.Priority,
	}
	if err := s.repo.CreateTarget(t); err != nil {
		return nil, err
	}

	// Reload with relations
	targets, _ := s.repo.FindTargetsByUser(userID)
	for _, tgt := range targets {
		if tgt.ID == t.ID {
			resp := ToUserTargetResponse(tgt)
			return &resp, nil
		}
	}
	resp := ToUserTargetResponse(*t)
	return &resp, nil
}

func (s *universityService) DeleteTarget(targetID uint, userID uint) error {
	targets, err := s.repo.FindTargetsByUser(userID)
	if err != nil {
		return err
	}
	for _, t := range targets {
		if t.ID == targetID {
			return s.repo.DeleteTarget(targetID)
		}
	}
	return errors.New("target not found or not owned by user")
}

func (s *universityService) GetUsersByUniversity(universityID uint) ([]UserResponse, error) {
	users, err := s.repo.FindUsersByUniversity(universityID)
	if err != nil {
		return nil, err
	}
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(u)
	}
	return result, nil
}
