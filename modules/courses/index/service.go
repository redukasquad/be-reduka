package courses

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/modules/dto"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type courseService struct {
	repo Repository
}

type Service interface {
	GetAll(params dto.ListQueryParams, requestID string) (*dto.PaginatedResponse[dto.CourseResponse], error)
	GetByID(id uint, requestID string) (*dto.CourseResponse, error)
	GetByProgramID(programID uint, requestID string) ([]dto.CourseResponse, error)
	Create(input CreateCourseInput, requestID string, userID uint) (*dto.CourseResponse, error)
	Update(id uint, input UpdateCourseInput, requestID string, userID uint) (*dto.CourseResponse, error)
	Delete(id uint, requestID string, userID uint) error
}

func NewService(repo Repository) Service {
	return &courseService{repo: repo}
}

func (s *courseService) GetAll(params dto.ListQueryParams, requestID string) (*dto.PaginatedResponse[dto.CourseResponse], error) {
	params.SetDefaults()

	utils.LogInfo("courses", "get_all", "Fetching courses with pagination", requestID, 0, map[string]any{
		"page":    params.Page,
		"perPage": params.PerPage,
		"search":  params.Q,
	})

	courses, err := s.repo.FindAllPaginated(params.GetOffset(), params.PerPage, params.Q)
	if err != nil {
		utils.LogError("courses", "get_all", "Failed to fetch courses: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	totalCount, err := s.repo.CountWithSearch(params.Q)
	if err != nil {
		utils.LogError("courses", "get_all", "Failed to count courses:"+err.Error(), requestID, 0, nil)
		return nil, err
	}

	var courseResponses []dto.CourseResponse

	for _, course := range courses {
		courseResponses = append(courseResponses, dto.ToCourseResponse(course))
	}

	response := dto.NewPaginatedResponse(courseResponses, params.Page, params.PerPage, totalCount)

	utils.LogSuccess("courses", "get_all", "Successfully fetched all courses", requestID, 0, map[string]any{
		"count":      len(courses),
		"totalItems": totalCount,
	})
	return &response, nil
}

func (s *courseService) GetByID(id uint, requestID string) (*dto.CourseResponse, error) {
	utils.LogInfo("courses", "get_by_id", "Fetching course by ID", requestID, 0, map[string]any{
		"course_id": id,
	})

	course, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("courses", "get_by_id", "Course not found", requestID, 0, map[string]any{
				"course_id": id,
			})
			return nil, errors.New("course not found")
		}
		utils.LogError("courses", "get_by_id", "Failed to fetch course: "+err.Error(), requestID, 0, map[string]any{
			"course_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("courses", "get_by_id", "Successfully fetched course", requestID, 0, map[string]any{
		"course_id":   course.ID,
		"course_name": course.NameCourse,
	})

	response := dto.ToCourseResponse(course)
	return &response, nil
}

func (s *courseService) GetByProgramID(programID uint, requestID string) ([]dto.CourseResponse, error) {
	utils.LogInfo("courses", "get_by_program_id", "Fetching courses by program ID", requestID, 0, map[string]any{
		"program_id": programID,
	})

	courses, err := s.repo.FindByProgramID(programID)
	if err != nil {
		utils.LogError("courses", "get_by_program_id", "Failed to fetch courses: "+err.Error(), requestID, 0, map[string]any{
			"program_id": programID,
		})
		return nil, err
	}

	var courseResponses []dto.CourseResponse
	for _, course := range courses {
		courseResponses = append(courseResponses, dto.ToCourseResponse(course))
	}

	utils.LogSuccess("courses", "get_by_program_id", "Successfully fetched courses by program", requestID, 0, map[string]any{
		"program_id": programID,
		"count":      len(courses),
	})
	return courseResponses, nil
}

func (s *courseService) Create(input CreateCourseInput, requestID string, userID uint) (*dto.CourseResponse, error) {
	utils.LogInfo("courses", "create", "Attempting to create new course", requestID, userID, map[string]any{
		"course_name": input.NameCourse,
		"program_id":  input.ProgramID,
	})

	_, err := s.repo.FindByName(input.NameCourse)
	if err == nil {
		utils.LogWarning("courses", "create", "Course with this name already exists", requestID, userID, map[string]any{
			"course_name": input.NameCourse,
		})
		return nil, errors.New("course with this name already exists")
	}

	course := &entities.Course{
		ProgramID:         input.ProgramID,
		CreatedByUserID:   userID,
		NameCourse:        input.NameCourse,
		Description:       input.Description,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		IsFree:            input.IsFree,
		WhatsappGroupLink: input.WhatsappGroupLink,
		Image: input.Image,
	}

	if err := s.repo.Create(course); err != nil {
		utils.LogError("courses", "create", "Failed to create course: "+err.Error(), requestID, userID, map[string]any{
			"course_name": input.NameCourse,
		})
		return nil, err
	}

	// Fetch dengan preload relations
	createdCourse, err := s.repo.FindByID(course.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("courses", "create", "Course created successfully", requestID, userID, map[string]any{
		"course_id":   course.ID,
		"course_name": course.NameCourse,
	})

	response := dto.ToCourseResponse(createdCourse)
	return &response, nil
}

func (s *courseService) Update(id uint, input UpdateCourseInput, requestID string, userID uint) (*dto.CourseResponse, error) {
	utils.LogInfo("courses", "update", "Attempting to update course", requestID, userID, map[string]any{
		"course_id": id,
	})

	course, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("courses", "update", "Course not found", requestID, userID, map[string]any{
				"course_id": id,
			})
			return nil, errors.New("course not found")
		}
		utils.LogError("courses", "update", "Failed to fetch course: "+err.Error(), requestID, userID, map[string]any{
			"course_id": id,
		})
		return nil, err
	}

	if input.ProgramID != nil {
		course.ProgramID = *input.ProgramID
	}
	if input.NameCourse != nil {
		if *input.NameCourse != course.NameCourse {
			existing, _ := s.repo.FindByName(*input.NameCourse)
			if existing.ID != 0 {
				utils.LogWarning("courses", "update", "Course with this name already exists", requestID, userID, map[string]any{
					"course_name": *input.NameCourse,
				})
				return nil, errors.New("course with this name already exists")
			}
		}
		course.NameCourse = *input.NameCourse
	}
	if input.Description != nil {
		course.Description = *input.Description
	}
	if input.StartDate != nil {
		course.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		course.EndDate = *input.EndDate
	}
	if input.IsFree != nil {
		course.IsFree = *input.IsFree
	}
	if input.WhatsappGroupLink != nil {
		course.WhatsappGroupLink = *input.WhatsappGroupLink
	}
	if input.Image != nil {
		course.Image = *input.Image
	}
	if err := s.repo.Update(&course); err != nil {
		utils.LogError("courses", "update", "Failed to update course: "+err.Error(), requestID, userID, map[string]any{
			"course_id": id,
		})
		return nil, err
	}

	// Fetch ulang dengan preload relations
	updatedCourse, err := s.repo.FindByID(course.ID)
	if err != nil {
		return nil, err
	}

	utils.LogSuccess("courses", "update", "Course updated successfully", requestID, userID, map[string]any{
		"course_id":   course.ID,
		"course_name": course.NameCourse,
	})

	response := dto.ToCourseResponse(updatedCourse)
	return &response, nil
}

func (s *courseService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("courses", "delete", "Attempting to delete course", requestID, userID, map[string]any{
		"course_id": id,
	})

	course, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("courses", "delete", "Course not found", requestID, userID, map[string]any{
				"course_id": id,
			})
			return errors.New("course not found")
		}
		utils.LogError("courses", "delete", "Failed to fetch course: "+err.Error(), requestID, userID, map[string]any{
			"course_id": id,
		})
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("courses", "delete", "Failed to delete course: "+err.Error(), requestID, userID, map[string]any{
			"course_id": id,
		})
		return err
	}

	utils.LogSuccess("courses", "delete", "Course deleted successfully", requestID, userID, map[string]any{
		"course_id":   id,
		"course_name": course.NameCourse,
	})
	return nil
}
