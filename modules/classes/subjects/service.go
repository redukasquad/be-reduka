package subjects

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type subjectService struct {
	repo Repository
}

type Service interface {
	GetByCourseID(courseID uint, requestID string) ([]SubjectResponse, error)
	GetByID(id uint, requestID string) (*SubjectResponse, error)
	Create(courseID uint, input CreateSubjectInput, requestID string, userID uint) (*SubjectResponse, error)
	Update(id uint, input UpdateSubjectInput, requestID string, userID uint) (*SubjectResponse, error)
	Delete(id uint, requestID string, userID uint) error
}

func NewService(repo Repository) Service {
	return &subjectService{repo: repo}
}

func (s *subjectService) GetByCourseID(courseID uint, requestID string) ([]SubjectResponse, error) {
	utils.LogInfo("classes", "get_by_course", "Fetching classes for course", requestID, 0, map[string]any{
		"course_id": courseID,
	})

	classes, err := s.repo.FindByCourseID(courseID)
	if err != nil {
		utils.LogError("classes", "get_by_course", "Failed to fetch classes: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []SubjectResponse
	for _, c := range classes {
		responses = append(responses, s.toResponse(c))
	}

	utils.LogSuccess("classes", "get_by_course", "Successfully fetched classes", requestID, 0, map[string]any{
		"course_id": courseID,
		"count":     len(responses),
	})
	return responses, nil
}

func (s *subjectService) GetByID(id uint, requestID string) (*SubjectResponse, error) {
	utils.LogInfo("classes", "get_by_id", "Fetching class by ID", requestID, 0, map[string]any{
		"class_id": id,
	})

	class, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		utils.LogError("classes", "get_by_id", "Failed to fetch class: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	response := s.toResponse(class)
	return &response, nil
}

func (s *subjectService) Create(courseID uint, input CreateSubjectInput, requestID string, userID uint) (*SubjectResponse, error) {
	utils.LogInfo("classes", "create", "Creating new class", requestID, userID, map[string]any{
		"course_id": courseID,
		"name":      input.Name,
	})

	class := &entities.Class{
		CourseID:    courseID,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.repo.Create(class); err != nil {
		utils.LogError("classes", "create", "Failed to create class: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("classes", "create", "Class created successfully", requestID, userID, map[string]any{
		"class_id": class.ID,
		"name":     class.Name,
	})

	response := s.toResponse(*class)
	return &response, nil
}

func (s *subjectService) Update(id uint, input UpdateSubjectInput, requestID string, userID uint) (*SubjectResponse, error) {
	utils.LogInfo("classes", "update", "Updating class", requestID, userID, map[string]any{
		"class_id": id,
	})

	class, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, err
	}

	if input.Name != nil {
		class.Name = *input.Name
	}
	if input.Description != nil {
		class.Description = *input.Description
	}

	if err := s.repo.Update(&class); err != nil {
		utils.LogError("classes", "update", "Failed to update class: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("classes", "update", "Class updated successfully", requestID, userID, map[string]any{
		"class_id": id,
	})

	response := s.toResponse(class)
	return &response, nil
}

func (s *subjectService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("classes", "delete", "Deleting class", requestID, userID, map[string]any{
		"class_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("classes", "delete", "Failed to delete class: "+err.Error(), requestID, userID, nil)
		return err
	}

	utils.LogSuccess("classes", "delete", "Class deleted successfully", requestID, userID, map[string]any{
		"class_id": id,
	})
	return nil
}

func (s *subjectService) toResponse(c entities.Class) SubjectResponse {
	return SubjectResponse{
		ID:          c.ID,
		CourseID:    c.CourseID,
		Name:        c.Name,
		Description: c.Description,
		LessonCount: len(c.Lessons),
	}
}
