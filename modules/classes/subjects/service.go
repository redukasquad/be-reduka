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
	utils.LogInfo("subjects", "get_by_course", "Fetching subjects for course", requestID, 0, map[string]any{
		"course_id": courseID,
	})

	subjects, err := s.repo.FindByCourseID(courseID)
	if err != nil {
		utils.LogError("subjects", "get_by_course", "Failed to fetch subjects: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []SubjectResponse
	for _, sub := range subjects {
		responses = append(responses, s.toSubjectResponse(sub))
	}

	utils.LogSuccess("subjects", "get_by_course", "Successfully fetched subjects", requestID, 0, map[string]any{
		"course_id": courseID,
		"count":     len(responses),
	})
	return responses, nil
}

func (s *subjectService) GetByID(id uint, requestID string) (*SubjectResponse, error) {
	utils.LogInfo("subjects", "get_by_id", "Fetching subject by ID", requestID, 0, map[string]any{
		"subject_id": id,
	})

	subject, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("subjects", "get_by_id", "Subject not found", requestID, 0, map[string]any{
				"subject_id": id,
			})
			return nil, errors.New("subject not found")
		}
		utils.LogError("subjects", "get_by_id", "Failed to fetch subject: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	response := s.toSubjectResponse(subject)
	utils.LogSuccess("subjects", "get_by_id", "Successfully fetched subject", requestID, 0, map[string]any{
		"subject_id": id,
	})
	return &response, nil
}

func (s *subjectService) Create(courseID uint, input CreateSubjectInput, requestID string, userID uint) (*SubjectResponse, error) {
	utils.LogInfo("subjects", "create", "Creating new subject", requestID, userID, map[string]any{
		"course_id": courseID,
		"name":      input.Name,
	})

	subject := &entities.ClassSubject{
		CourseID:    courseID,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := s.repo.Create(subject); err != nil {
		utils.LogError("subjects", "create", "Failed to create subject: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("subjects", "create", "Subject created successfully", requestID, userID, map[string]any{
		"subject_id": subject.ID,
		"name":       subject.Name,
	})

	response := s.toSubjectResponse(*subject)
	return &response, nil
}

func (s *subjectService) Update(id uint, input UpdateSubjectInput, requestID string, userID uint) (*SubjectResponse, error) {
	utils.LogInfo("subjects", "update", "Updating subject", requestID, userID, map[string]any{
		"subject_id": id,
	})

	subject, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("subjects", "update", "Subject not found", requestID, userID, map[string]any{
				"subject_id": id,
			})
			return nil, errors.New("subject not found")
		}
		return nil, err
	}

	if input.Name != nil {
		subject.Name = *input.Name
	}
	if input.Description != nil {
		subject.Description = *input.Description
	}

	if err := s.repo.Update(&subject); err != nil {
		utils.LogError("subjects", "update", "Failed to update subject: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("subjects", "update", "Subject updated successfully", requestID, userID, map[string]any{
		"subject_id": id,
	})

	response := s.toSubjectResponse(subject)
	return &response, nil
}

func (s *subjectService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("subjects", "delete", "Deleting subject", requestID, userID, map[string]any{
		"subject_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("subjects", "delete", "Subject not found", requestID, userID, map[string]any{
				"subject_id": id,
			})
			return errors.New("subject not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("subjects", "delete", "Failed to delete subject: "+err.Error(), requestID, userID, nil)
		return err
	}

	utils.LogSuccess("subjects", "delete", "Subject deleted successfully", requestID, userID, map[string]any{
		"subject_id": id,
	})
	return nil
}

func (s *subjectService) toSubjectResponse(sub entities.ClassSubject) SubjectResponse {
	return SubjectResponse{
		ID:          sub.ID,
		CourseID:    sub.CourseID,
		Name:        sub.Name,
		Description: sub.Description,
		LessonCount: len(sub.Lessons),
	}
}
