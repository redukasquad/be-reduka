package lessons

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type lessonService struct {
	repo Repository
}

type Service interface {
	GetByClassID(classID uint, requestID string) ([]LessonResponse, error)
	GetByID(id uint, requestID string) (*LessonResponse, error)
	Create(classID uint, input CreateLessonInput, requestID string, userID uint) (*LessonResponse, error)
	Update(id uint, input UpdateLessonInput, requestID string, userID uint) (*LessonResponse, error)
	Delete(id uint, requestID string, userID uint) error
}

func NewService(repo Repository) Service {
	return &lessonService{repo: repo}
}

func (s *lessonService) GetByClassID(classID uint, requestID string) ([]LessonResponse, error) {
	utils.LogInfo("lessons", "get_by_class", "Fetching lessons for class", requestID, 0, map[string]any{
		"class_id": classID,
	})

	lessons, err := s.repo.FindByClassID(classID)
	if err != nil {
		utils.LogError("lessons", "get_by_class", "Failed to fetch lessons: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []LessonResponse
	for _, lesson := range lessons {
		responses = append(responses, s.toResponse(lesson, true))
	}

	utils.LogSuccess("lessons", "get_by_class", "Successfully fetched lessons", requestID, 0, map[string]any{
		"class_id": classID,
		"count":    len(responses),
	})
	return responses, nil
}

func (s *lessonService) GetByID(id uint, requestID string) (*LessonResponse, error) {
	utils.LogInfo("lessons", "get_by_id", "Fetching lesson by ID", requestID, 0, map[string]any{
		"lesson_id": id,
	})

	lesson, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		utils.LogError("lessons", "get_by_id", "Failed to fetch lesson: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	response := s.toResponse(lesson, true)
	return &response, nil
}

func (s *lessonService) Create(classID uint, input CreateLessonInput, requestID string, userID uint) (*LessonResponse, error) {
	utils.LogInfo("lessons", "create", "Creating new lesson", requestID, userID, map[string]any{
		"class_id": classID,
		"title":    input.Title,
	})

	lesson := &entities.Lesson{
		ClassID:         classID,
		CreatedByUserID: userID,
		Title:           input.Title,
		Description:     input.Description,
		LessonOrder:     input.LessonOrder,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
	}

	if err := s.repo.Create(lesson); err != nil {
		utils.LogError("lessons", "create", "Failed to create lesson: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("lessons", "create", "Lesson created successfully", requestID, userID, map[string]any{
		"lesson_id": lesson.ID,
		"title":     lesson.Title,
	})

	response := s.toResponse(*lesson, false)
	return &response, nil
}

func (s *lessonService) Update(id uint, input UpdateLessonInput, requestID string, userID uint) (*LessonResponse, error) {
	utils.LogInfo("lessons", "update", "Updating lesson", requestID, userID, map[string]any{
		"lesson_id": id,
	})

	lesson, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}

	if input.Title != nil {
		lesson.Title = *input.Title
	}
	if input.Description != nil {
		lesson.Description = *input.Description
	}
	if input.LessonOrder != nil {
		lesson.LessonOrder = *input.LessonOrder
	}
	if input.StartTime != nil {
		lesson.StartTime = input.StartTime
	}
	if input.EndTime != nil {
		lesson.EndTime = input.EndTime
	}

	if err := s.repo.Update(&lesson); err != nil {
		utils.LogError("lessons", "update", "Failed to update lesson: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("lessons", "update", "Lesson updated successfully", requestID, userID, map[string]any{
		"lesson_id": id,
	})

	response := s.toResponse(lesson, true)
	return &response, nil
}

func (s *lessonService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("lessons", "delete", "Deleting lesson", requestID, userID, map[string]any{
		"lesson_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("lesson not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("lessons", "delete", "Failed to delete lesson: "+err.Error(), requestID, userID, nil)
		return err
	}

	utils.LogSuccess("lessons", "delete", "Lesson deleted successfully", requestID, userID, map[string]any{
		"lesson_id": id,
	})
	return nil
}

func (s *lessonService) toResponse(lesson entities.Lesson, includeResources bool) LessonResponse {
	response := LessonResponse{
		ID:            lesson.ID,
		ClassID:       lesson.ClassID,
		Title:         lesson.Title,
		Description:   lesson.Description,
		LessonOrder:   lesson.LessonOrder,
		StartTime:     lesson.StartTime,
		EndTime:       lesson.EndTime,
		ResourceCount: len(lesson.Resources),
	}

	if lesson.Class.ID != 0 {
		response.ClassName = lesson.Class.Name
	}

	if includeResources && len(lesson.Resources) > 0 {
		for _, res := range lesson.Resources {
			response.Resources = append(response.Resources, ResourceResponse{
				ID:    res.ID,
				Type:  res.Type,
				Title: res.Title,
				URL:   res.URL,
			})
		}
	}

	return response
}
