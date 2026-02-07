package tryouts

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/modules/dto"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type tryOutService struct {
	repo Repository
}

type Service interface {
	// Try Out
	GetAll(params dto.ListQueryParams, requestID string, isAdmin bool) (*dto.PaginatedResponse[TryOutBriefResponse], error)
	GetByID(id uint, requestID string) (*TryOutResponse, error)
	Create(input CreateTryOutInput, requestID string, userID uint) (*TryOutResponse, error)
	Update(id uint, input UpdateTryOutInput, requestID string, userID uint) (*TryOutResponse, error)
	Delete(id uint, requestID string, userID uint) error

	// Tutor Permissions
	GetTutorPermissions(tryOutID uint, requestID string) ([]TutorPermissionResponse, error)
	GrantTutorPermission(tryOutID uint, input GrantTutorPermissionInput, requestID string, grantedByUserID uint) (*TutorPermissionResponse, error)
	RevokeTutorPermission(tryOutID, userID uint, requestID string, revokedByUserID uint) error
	HasTutorPermission(tryOutID, userID uint) (bool, error)
}

func NewService(repo Repository) Service {
	return &tryOutService{repo: repo}
}

// ==========================================
// Try Out Service Methods
// ==========================================

func (s *tryOutService) GetAll(params dto.ListQueryParams, requestID string, isAdmin bool) (*dto.PaginatedResponse[TryOutBriefResponse], error) {
	params.SetDefaults()

	utils.LogInfo("tryouts", "get_all", "Fetching try outs with pagination", requestID, 0, map[string]any{
		"page":    params.Page,
		"perPage": params.PerPage,
		"search":  params.Q,
		"isAdmin": isAdmin,
	})

	// Non-admin users only see published try outs
	publishedOnly := !isAdmin

	tryOuts, err := s.repo.FindAllPaginated(params.GetOffset(), params.PerPage, params.Q, publishedOnly)
	if err != nil {
		utils.LogError("tryouts", "get_all", "Failed to fetch try outs: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	totalCount, err := s.repo.CountWithSearch(params.Q, publishedOnly)
	if err != nil {
		utils.LogError("tryouts", "get_all", "Failed to count try outs: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []TryOutBriefResponse
	for _, tryOut := range tryOuts {
		responses = append(responses, toTryOutBriefResponse(tryOut))
	}

	response := dto.NewPaginatedResponse(responses, params.Page, params.PerPage, totalCount)

	utils.LogSuccess("tryouts", "get_all", "Successfully fetched try outs", requestID, 0, map[string]any{
		"count":      len(tryOuts),
		"totalItems": totalCount,
	})
	return &response, nil
}

func (s *tryOutService) GetByID(id uint, requestID string) (*TryOutResponse, error) {
	utils.LogInfo("tryouts", "get_by_id", "Fetching try out by ID", requestID, 0, map[string]any{
		"try_out_id": id,
	})

	tryOut, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("tryouts", "get_by_id", "Try out not found", requestID, 0, map[string]any{
				"try_out_id": id,
			})
			return nil, errors.New("try out not found")
		}
		utils.LogError("tryouts", "get_by_id", "Failed to fetch try out: "+err.Error(), requestID, 0, map[string]any{
			"try_out_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("tryouts", "get_by_id", "Successfully fetched try out", requestID, 0, map[string]any{
		"try_out_id":   tryOut.ID,
		"try_out_name": tryOut.Name,
	})

	response := toTryOutResponse(tryOut)
	return &response, nil
}

func (s *tryOutService) Create(input CreateTryOutInput, requestID string, userID uint) (*TryOutResponse, error) {
	utils.LogInfo("tryouts", "create", "Attempting to create new try out", requestID, userID, map[string]any{
		"try_out_name": input.Name,
	})

	// Check if name already exists
	_, err := s.repo.FindByName(input.Name)
	if err == nil {
		utils.LogWarning("tryouts", "create", "Try out with this name already exists", requestID, userID, map[string]any{
			"try_out_name": input.Name,
		})
		return nil, errors.New("try out with this name already exists")
	}

	// Validate registration dates
	if input.RegistrationEnd.Before(input.RegistrationStart) {
		return nil, errors.New("registration end date must be after start date")
	}

	tryOut := &entities.TryOut{
		Name:              input.Name,
		Description:       input.Description,
		ImageURL:          input.ImageURL,
		IsFree:            input.IsFree,
		Price:             input.Price,
		QrisImageURL:      input.QrisImageURL,
		PaymentLink:       input.PaymentLink,
		RegistrationStart: input.RegistrationStart,
		RegistrationEnd:   input.RegistrationEnd,
		IsPublished:       input.IsPublished,
		CreatedByUserID:   userID,
	}

	if err := s.repo.Create(tryOut); err != nil {
		utils.LogError("tryouts", "create", "Failed to create try out: "+err.Error(), requestID, userID, map[string]any{
			"try_out_name": input.Name,
		})
		return nil, err
	}

	// Fetch with preload
	createdTryOut, err := s.repo.FindByID(tryOut.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("tryouts", "create", "Try out created successfully", requestID, userID, map[string]any{
		"try_out_id":   tryOut.ID,
		"try_out_name": tryOut.Name,
	})

	response := toTryOutResponse(createdTryOut)
	return &response, nil
}

func (s *tryOutService) Update(id uint, input UpdateTryOutInput, requestID string, userID uint) (*TryOutResponse, error) {
	utils.LogInfo("tryouts", "update", "Attempting to update try out", requestID, userID, map[string]any{
		"try_out_id": id,
	})

	tryOut, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("tryouts", "update", "Try out not found", requestID, userID, map[string]any{
				"try_out_id": id,
			})
			return nil, errors.New("try out not found")
		}
		utils.LogError("tryouts", "update", "Failed to fetch try out: "+err.Error(), requestID, userID, map[string]any{
			"try_out_id": id,
		})
		return nil, err
	}

	// Update fields if provided
	if input.Name != nil {
		if *input.Name != tryOut.Name {
			existing, _ := s.repo.FindByName(*input.Name)
			if existing.ID != 0 {
				return nil, errors.New("try out with this name already exists")
			}
		}
		tryOut.Name = *input.Name
	}
	if input.Description != nil {
		tryOut.Description = *input.Description
	}
	if input.ImageURL != nil {
		tryOut.ImageURL = *input.ImageURL
	}
	if input.IsFree != nil {
		tryOut.IsFree = *input.IsFree
	}
	if input.Price != nil {
		tryOut.Price = *input.Price
	}
	if input.QrisImageURL != nil {
		tryOut.QrisImageURL = *input.QrisImageURL
	}
	if input.PaymentLink != nil {
		tryOut.PaymentLink = *input.PaymentLink
	}
	if input.RegistrationStart != nil {
		tryOut.RegistrationStart = *input.RegistrationStart
	}
	if input.RegistrationEnd != nil {
		tryOut.RegistrationEnd = *input.RegistrationEnd
	}
	if input.IsPublished != nil {
		tryOut.IsPublished = *input.IsPublished
	}

	// Validate registration dates
	if tryOut.RegistrationEnd.Before(tryOut.RegistrationStart) {
		return nil, errors.New("registration end date must be after start date")
	}

	if err := s.repo.Update(&tryOut); err != nil {
		utils.LogError("tryouts", "update", "Failed to update try out: "+err.Error(), requestID, userID, map[string]any{
			"try_out_id": id,
		})
		return nil, err
	}

	// Fetch with preload
	updatedTryOut, err := s.repo.FindByID(tryOut.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("tryouts", "update", "Try out updated successfully", requestID, userID, map[string]any{
		"try_out_id":   tryOut.ID,
		"try_out_name": tryOut.Name,
	})

	response := toTryOutResponse(updatedTryOut)
	return &response, nil
}

func (s *tryOutService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("tryouts", "delete", "Attempting to delete try out", requestID, userID, map[string]any{
		"try_out_id": id,
	})

	tryOut, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("tryouts", "delete", "Try out not found", requestID, userID, map[string]any{
				"try_out_id": id,
			})
			return errors.New("try out not found")
		}
		utils.LogError("tryouts", "delete", "Failed to fetch try out: "+err.Error(), requestID, userID, map[string]any{
			"try_out_id": id,
		})
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("tryouts", "delete", "Failed to delete try out: "+err.Error(), requestID, userID, map[string]any{
			"try_out_id": id,
		})
		return err
	}

	utils.LogSuccess("tryouts", "delete", "Try out deleted successfully", requestID, userID, map[string]any{
		"try_out_id":   id,
		"try_out_name": tryOut.Name,
	})
	return nil
}

// ==========================================
// Tutor Permission Service Methods
// ==========================================

func (s *tryOutService) GetTutorPermissions(tryOutID uint, requestID string) ([]TutorPermissionResponse, error) {
	utils.LogInfo("tryouts", "get_tutor_permissions", "Fetching tutor permissions", requestID, 0, map[string]any{
		"try_out_id": tryOutID,
	})

	permissions, err := s.repo.FindTutorPermissions(tryOutID)
	if err != nil {
		utils.LogError("tryouts", "get_tutor_permissions", "Failed to fetch permissions: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var responses []TutorPermissionResponse
	for _, p := range permissions {
		responses = append(responses, toTutorPermissionResponse(p))
	}

	return responses, nil
}

func (s *tryOutService) GrantTutorPermission(tryOutID uint, input GrantTutorPermissionInput, requestID string, grantedByUserID uint) (*TutorPermissionResponse, error) {
	utils.LogInfo("tryouts", "grant_tutor_permission", "Granting tutor permission", requestID, grantedByUserID, map[string]any{
		"try_out_id": tryOutID,
		"user_id":    input.UserID,
	})

	// Check if try out exists
	_, err := s.repo.FindByID(tryOutID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("try out not found")
		}
		return nil, err
	}

	// Check if permission already exists
	_, err = s.repo.FindTutorPermission(tryOutID, input.UserID)
	if err == nil {
		return nil, errors.New("tutor already has permission for this try out")
	}

	permission := &entities.TutorPermission{
		TryOutPackageID: tryOutID,
		UserID:          input.UserID,
		GrantedByUserID: grantedByUserID,
	}

	if err := s.repo.CreateTutorPermission(permission); err != nil {
		utils.LogError("tryouts", "grant_tutor_permission", "Failed to grant permission: "+err.Error(), requestID, grantedByUserID, nil)
		return nil, err
	}

	// Fetch with preload
	createdPermission, err := s.repo.FindTutorPermission(tryOutID, input.UserID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("tryouts", "grant_tutor_permission", "Permission granted successfully", requestID, grantedByUserID, map[string]any{
		"try_out_id": tryOutID,
		"user_id":    input.UserID,
	})

	response := toTutorPermissionResponse(createdPermission)
	return &response, nil
}

func (s *tryOutService) RevokeTutorPermission(tryOutID, userID uint, requestID string, revokedByUserID uint) error {
	utils.LogInfo("tryouts", "revoke_tutor_permission", "Revoking tutor permission", requestID, revokedByUserID, map[string]any{
		"try_out_id": tryOutID,
		"user_id":    userID,
	})

	// Check if permission exists
	_, err := s.repo.FindTutorPermission(tryOutID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tutor permission not found")
		}
		return err
	}

	if err := s.repo.DeleteTutorPermission(tryOutID, userID); err != nil {
		utils.LogError("tryouts", "revoke_tutor_permission", "Failed to revoke permission: "+err.Error(), requestID, revokedByUserID, nil)
		return err
	}

	utils.LogSuccess("tryouts", "revoke_tutor_permission", "Permission revoked successfully", requestID, revokedByUserID, map[string]any{
		"try_out_id": tryOutID,
		"user_id":    userID,
	})
	return nil
}

func (s *tryOutService) HasTutorPermission(tryOutID, userID uint) (bool, error) {
	return s.repo.HasTutorPermission(tryOutID, userID)
}

// ==========================================
// Helper Functions - Entity to Response Mapping
// ==========================================

func toTryOutResponse(tryOut entities.TryOut) TryOutResponse {
	response := TryOutResponse{
		ID:                tryOut.ID,
		Name:              tryOut.Name,
		Description:       tryOut.Description,
		ImageURL:          tryOut.ImageURL,
		IsFree:            tryOut.IsFree,
		Price:             tryOut.Price,
		QrisImageURL:      tryOut.QrisImageURL,
		PaymentLink:       tryOut.PaymentLink,
		RegistrationStart: tryOut.RegistrationStart,
		RegistrationEnd:   tryOut.RegistrationEnd,
		IsPublished:       tryOut.IsPublished,
		CreatedAt:         tryOut.CreatedAt,
	}

	if tryOut.Creator.ID != 0 {
		creator := toCreatorBriefResponse(tryOut.Creator)
		response.Creator = &creator
	}

	return response
}

func toTryOutBriefResponse(tryOut entities.TryOut) TryOutBriefResponse {
	return TryOutBriefResponse{
		ID:                tryOut.ID,
		Name:              tryOut.Name,
		ImageURL:          tryOut.ImageURL,
		IsFree:            tryOut.IsFree,
		Price:             tryOut.Price,
		RegistrationStart: tryOut.RegistrationStart,
		RegistrationEnd:   tryOut.RegistrationEnd,
		IsPublished:       tryOut.IsPublished,
	}
}

func toCreatorBriefResponse(user entities.User) CreatorBriefResponse {
	return CreatorBriefResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func toTutorPermissionResponse(permission entities.TutorPermission) TutorPermissionResponse {
	response := TutorPermissionResponse{
		ID:        permission.ID,
		TryOutID:  permission.TryOutPackageID,
		GrantedAt: permission.GrantedAt,
	}

	if permission.User.ID != 0 {
		user := toCreatorBriefResponse(permission.User)
		response.User = &user
	}

	if permission.GrantedBy.ID != 0 {
		grantedBy := toCreatorBriefResponse(permission.GrantedBy)
		response.GrantedBy = &grantedBy
	}

	return response
}
