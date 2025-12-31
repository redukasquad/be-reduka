package resources

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type resourceService struct {
	repo Repository
}

type Service interface {
	GetByLessonID(lessonID uint, requestID string) ([]ResourceResponse, error)
	Create(lessonID uint, input CreateResourceInput, requestID string, userID uint) (*ResourceResponse, error)
	Update(id uint, input UpdateResourceInput, requestID string, userID uint) (*ResourceResponse, error)
	Delete(id uint, requestID string, userID uint) error
}

func NewService(repo Repository) Service {
	return &resourceService{repo: repo}
}

func (s *resourceService) GetByLessonID(lessonID uint, requestID string) ([]ResourceResponse, error) {
	utils.LogInfo("resources", "get_by_lesson", "Fetching resources for lesson", requestID, 0, map[string]any{
		"lesson_id": lessonID,
	})

	resources, err := s.repo.FindByLessonID(lessonID)
	if err != nil {
		utils.LogError("resources", "get_by_lesson", "Failed to fetch resources: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []ResourceResponse
	for _, res := range resources {
		responses = append(responses, s.toResourceResponse(res))
	}

	utils.LogSuccess("resources", "get_by_lesson", "Successfully fetched resources", requestID, 0, map[string]any{
		"lesson_id": lessonID,
		"count":     len(responses),
	})
	return responses, nil
}

func (s *resourceService) Create(lessonID uint, input CreateResourceInput, requestID string, userID uint) (*ResourceResponse, error) {
	utils.LogInfo("resources", "create", "Creating new resource", requestID, userID, map[string]any{
		"lesson_id": lessonID,
		"type":      input.Type,
	})

	resource := &entities.ClassLessonResource{
		ClassLessonID: lessonID,
		Type:          input.Type,
		Title:         input.Title,
		URL:           input.URL,
	}

	if err := s.repo.Create(resource); err != nil {
		utils.LogError("resources", "create", "Failed to create resource: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("resources", "create", "Resource created successfully", requestID, userID, map[string]any{
		"resource_id": resource.ID,
		"type":        resource.Type,
	})

	response := s.toResourceResponse(*resource)
	return &response, nil
}

func (s *resourceService) Update(id uint, input UpdateResourceInput, requestID string, userID uint) (*ResourceResponse, error) {
	utils.LogInfo("resources", "update", "Updating resource", requestID, userID, map[string]any{
		"resource_id": id,
	})

	resource, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("resources", "update", "Resource not found", requestID, userID, map[string]any{
				"resource_id": id,
			})
			return nil, errors.New("resource not found")
		}
		return nil, err
	}

	if input.Type != nil {
		resource.Type = *input.Type
	}
	if input.Title != nil {
		resource.Title = *input.Title
	}
	if input.URL != nil {
		resource.URL = *input.URL
	}

	if err := s.repo.Update(&resource); err != nil {
		utils.LogError("resources", "update", "Failed to update resource: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	utils.LogSuccess("resources", "update", "Resource updated successfully", requestID, userID, map[string]any{
		"resource_id": id,
	})

	response := s.toResourceResponse(resource)
	return &response, nil
}

func (s *resourceService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("resources", "delete", "Deleting resource", requestID, userID, map[string]any{
		"resource_id": id,
	})

	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("resources", "delete", "Resource not found", requestID, userID, map[string]any{
				"resource_id": id,
			})
			return errors.New("resource not found")
		}
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("resources", "delete", "Failed to delete resource: "+err.Error(), requestID, userID, nil)
		return err
	}

	utils.LogSuccess("resources", "delete", "Resource deleted successfully", requestID, userID, map[string]any{
		"resource_id": id,
	})
	return nil
}

func (s *resourceService) toResourceResponse(res entities.ClassLessonResource) ResourceResponse {
	response := ResourceResponse{
		ID:            res.ID,
		ClassLessonID: res.ClassLessonID,
		Type:          res.Type,
		Title:         res.Title,
		URL:           res.URL,
	}

	if res.ClassLesson.ID != 0 {
		response.LessonTitle = res.ClassLesson.Title
	}

	return response
}
