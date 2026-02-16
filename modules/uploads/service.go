package uploads

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type uploadService struct {
	repo Repository
}

type Service interface {
	Create(input CreateImageInput, requestID string) (*ImageResponse, error)
	DeleteByURL(url string, requestID string) (string, error)
}

func NewService(repo Repository) Service {
	return &uploadService{repo: repo}
}

func (s *uploadService) Create(input CreateImageInput, requestID string) (*ImageResponse, error) {
	utils.LogInfo("uploads", "create", "Attempting to create image record", requestID, 0, map[string]any{
		"url":    input.URL,
		"fileId": input.Fileid,
	})

	// Check if URL already exists
	_, err := s.repo.FindByURL(input.URL)
	if err == nil {
		utils.LogWarning("uploads", "create", "Image with this URL already exists", requestID, 0, map[string]any{
			"url": input.URL,
		})
		return nil, errors.New("image with this URL already exists")
	}

	image := &entities.Image{
		URL:    input.URL,
		Fileid: input.Fileid,
	}

	if err := s.repo.Create(image); err != nil {
		utils.LogError("uploads", "create", "Failed to create image: "+err.Error(), requestID, 0, map[string]any{
			"url": input.URL,
		})
		return nil, err
	}

	utils.LogSuccess("uploads", "create", "Image created successfully", requestID, 0, map[string]any{
		"image_id": image.ID,
		"url":      image.URL,
	})

	response := &ImageResponse{
		ID:        image.ID,
		URL:       image.URL,
		Fileid:    image.Fileid,
		CreatedAt: image.CreatedAt,
	}
	return response, nil
}

func (s *uploadService) DeleteByURL(url string, requestID string) (string, error) {
	utils.LogInfo("uploads", "delete", "Attempting to delete image by URL", requestID, 0, map[string]any{
		"url": url,
	})

	// Check if image exists
	image, err := s.repo.FindByURL(url)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("uploads", "delete", "Image not found", requestID, 0, map[string]any{
				"url": url,
			})
			return "", errors.New("image not found")
		}
		utils.LogError("uploads", "delete", "Failed to find image: "+err.Error(), requestID, 0, map[string]any{
			"url": url,
		})
		return "", err
	}

	if err := s.repo.DeleteByURL(url); err != nil {
		utils.LogError("uploads", "delete", "Failed to delete image: "+err.Error(), requestID, 0, map[string]any{
			"url": url,
		})
		return "", err
	}

	utils.LogSuccess("uploads", "delete", "Image deleted successfully", requestID, 0, map[string]any{
		"url": url,
	})
	return image.Fileid, nil
}
