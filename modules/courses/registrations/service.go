package registrations

import (
	"errors"
	"time"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type registrationService struct {
	repo Repository
}

type Service interface {
	Register(courseID uint, userID uint, input RegisterCourseInput, requestID string) (*RegistrationResponse, error)
	GetMyRegistrations(userID uint, requestID string) ([]RegistrationResponse, error)
	GetRegistrationsByCourse(courseID uint, requestID string) ([]RegistrationResponse, error)
	GetRegistrationByID(id uint, requestID string) (*RegistrationResponse, error)
	ApproveRegistration(id uint, requestID string, adminUserID uint) (*RegistrationResponse, error)
	RejectRegistration(id uint, requestID string, adminUserID uint) (*RegistrationResponse, error)
}

func NewService(repo Repository) Service {
	return &registrationService{repo: repo}
}

func (s *registrationService) Register(courseID uint, userID uint, input RegisterCourseInput, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "register", "Attempting to register user for course", requestID, userID, map[string]any{
		"course_id": courseID,
	})

	existing, err := s.repo.FindByUserAndCourse(userID, courseID)
	if err == nil && existing.ID != 0 {
		utils.LogWarning("registrations", "register", "User already registered for this course", requestID, userID, map[string]any{
			"course_id": courseID,
		})
		return nil, errors.New("you have already registered for this course")
	}

	registration := &entities.CourseRegistration{
		UserID:   userID,
		CourseID: courseID,
		Status:   "pending",
	}

	if err := s.repo.Create(registration); err != nil {
		utils.LogError("registrations", "register", "Failed to create registration: "+err.Error(), requestID, userID, map[string]any{
			"course_id": courseID,
		})
		return nil, err
	}

	if len(input.Answers) > 0 {
		var answers []entities.RegistrationAnswer
		for _, ans := range input.Answers {
			answers = append(answers, entities.RegistrationAnswer{
				RegistrationID: registration.ID,
				QuestionID:     ans.QuestionID,
				AnswerText:     ans.AnswerText,
			})
		}
		if err := s.repo.CreateAnswers(answers); err != nil {
			utils.LogError("registrations", "register", "Failed to create answers: "+err.Error(), requestID, userID, map[string]any{
				"course_id":       courseID,
				"registration_id": registration.ID,
			})
		}
	}

	utils.LogSuccess("registrations", "register", "Registration created successfully", requestID, userID, map[string]any{
		"course_id":       courseID,
		"registration_id": registration.ID,
		"status":          "pending",
	})

	fullReg, _ := s.repo.FindByID(registration.ID)
	return s.toRegistrationResponse(fullReg, false), nil
}

func (s *registrationService) GetMyRegistrations(userID uint, requestID string) ([]RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_my_registrations", "Fetching user registrations", requestID, userID, nil)

	registrations, err := s.repo.FindByUserID(userID)
	if err != nil {
		utils.LogError("registrations", "get_my_registrations", "Failed to fetch registrations: "+err.Error(), requestID, userID, nil)
		return nil, err
	}

	var responses []RegistrationResponse
	for _, reg := range registrations {
		showWhatsApp := reg.Status == "approved"
		responses = append(responses, *s.toRegistrationResponse(reg, showWhatsApp))
	}

	utils.LogSuccess("registrations", "get_my_registrations", "Successfully fetched user registrations", requestID, userID, map[string]any{
		"count": len(responses),
	})
	return responses, nil
}

func (s *registrationService) GetRegistrationsByCourse(courseID uint, requestID string) ([]RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_by_course", "Fetching registrations for course", requestID, 0, map[string]any{
		"course_id": courseID,
	})

	registrations, err := s.repo.FindByCourseID(courseID)
	if err != nil {
		utils.LogError("registrations", "get_by_course", "Failed to fetch registrations: "+err.Error(), requestID, 0, map[string]any{
			"course_id": courseID,
		})
		return nil, err
	}

	var responses []RegistrationResponse
	for _, reg := range registrations {
		responses = append(responses, *s.toRegistrationResponse(reg, true))
	}

	utils.LogSuccess("registrations", "get_by_course", "Successfully fetched course registrations", requestID, 0, map[string]any{
		"course_id": courseID,
		"count":     len(responses),
	})
	return responses, nil
}

func (s *registrationService) GetRegistrationByID(id uint, requestID string) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "get_by_id", "Fetching registration by ID", requestID, 0, map[string]any{
		"registration_id": id,
	})

	registration, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("registrations", "get_by_id", "Registration not found", requestID, 0, map[string]any{
				"registration_id": id,
			})
			return nil, errors.New("registration not found")
		}
		utils.LogError("registrations", "get_by_id", "Failed to fetch registration: "+err.Error(), requestID, 0, map[string]any{
			"registration_id": id,
		})
		return nil, err
	}

	return s.toRegistrationResponse(registration, true), nil
}

func (s *registrationService) ApproveRegistration(id uint, requestID string, adminUserID uint) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "approve", "Attempting to approve registration", requestID, adminUserID, map[string]any{
		"registration_id": id,
	})

	registration, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("registrations", "approve", "Registration not found", requestID, adminUserID, map[string]any{
				"registration_id": id,
			})
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	if registration.Status != "pending" {
		utils.LogWarning("registrations", "approve", "Registration is not pending", requestID, adminUserID, map[string]any{
			"registration_id": id,
			"current_status":  registration.Status,
		})
		return nil, errors.New("registration is not in pending status")
	}

	registration.Status = "approved"
	if err := s.repo.Update(&registration); err != nil {
		utils.LogError("registrations", "approve", "Failed to approve registration: "+err.Error(), requestID, adminUserID, map[string]any{
			"registration_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("registrations", "approve", "Registration approved successfully", requestID, adminUserID, map[string]any{
		"registration_id": id,
		"user_id":         registration.UserID,
		"course_id":       registration.CourseID,
	})

	registration, _ = s.repo.FindByID(id)
	return s.toRegistrationResponse(registration, true), nil
}

func (s *registrationService) RejectRegistration(id uint, requestID string, adminUserID uint) (*RegistrationResponse, error) {
	utils.LogInfo("registrations", "reject", "Attempting to reject registration", requestID, adminUserID, map[string]any{
		"registration_id": id,
	})

	registration, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("registrations", "reject", "Registration not found", requestID, adminUserID, map[string]any{
				"registration_id": id,
			})
			return nil, errors.New("registration not found")
		}
		return nil, err
	}

	if registration.Status != "pending" {
		utils.LogWarning("registrations", "reject", "Registration is not pending", requestID, adminUserID, map[string]any{
			"registration_id": id,
			"current_status":  registration.Status,
		})
		return nil, errors.New("registration is not in pending status")
	}

	registration.Status = "rejected"
	if err := s.repo.Update(&registration); err != nil {
		utils.LogError("registrations", "reject", "Failed to reject registration: "+err.Error(), requestID, adminUserID, map[string]any{
			"registration_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("registrations", "reject", "Registration rejected successfully", requestID, adminUserID, map[string]any{
		"registration_id": id,
		"user_id":         registration.UserID,
		"course_id":       registration.CourseID,
	})

	registration, _ = s.repo.FindByID(id)
	return s.toRegistrationResponse(registration, false), nil
}

func (s *registrationService) toRegistrationResponse(reg entities.CourseRegistration, showWhatsApp bool) *RegistrationResponse {
	response := &RegistrationResponse{
		ID:        reg.ID,
		UserID:    reg.UserID,
		CourseID:  reg.CourseID,
		Status:    reg.Status,
		CreatedAt: reg.CreatedAt.Format(time.RFC3339),
	}

	if reg.Course.ID != 0 {
		response.CourseName = reg.Course.NameCourse
		if showWhatsApp {
			response.WhatsappGroupLink = reg.Course.WhatsappGroupLink
		}
		if reg.Course.Program.ID != 0 {
			response.ProgramName = reg.Course.Program.ProgramName
		}
	}

	if reg.User.ID != 0 {
		response.UserName = reg.User.Username
		response.UserEmail = reg.User.Email
	}

	if len(reg.Answers) > 0 {
		for _, ans := range reg.Answers {
			response.Answers = append(response.Answers, AnswerResponse{
				QuestionID:   ans.QuestionID,
				QuestionText: ans.Question.QuestionText,
				AnswerText:   ans.AnswerText,
			})
		}
	}

	return response
}
