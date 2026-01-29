package programs

import (
	"errors"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/redukasquad/be-reduka/packages/utils"
	"gorm.io/gorm"
)

type programService struct {
	repo Repository
}

// Service interface defines the business logic methods for programs
type Service interface {
	GetAll(requestID string) ([]entities.Program, error)
	GetByID(id uint, requestID string) (*entities.Program, error)
	GetByName(name string, requestID string) (*entities.Program, error)
	Create(input CreateProgramInput, requestID string, userID uint) (*entities.Program, error)
	Update(id uint, input UpdateProgramInput, requestID string, userID uint) (*entities.Program, error)
	Delete(id uint, requestID string, userID uint) error
}

// NewService creates a new program service
func NewService(repo Repository) Service {
	return &programService{repo: repo}
}

func (s *programService) GetAll(requestID string) ([]entities.Program, error) {
	utils.LogInfo("programs", "get_all", "Fetching all programs", requestID, 0, nil)

	programs, err := s.repo.FindAll()
	if err != nil {
		utils.LogError("programs", "get_all", "Failed to fetch programs: "+err.Error(), requestID, 0, nil)
		return nil, err
	}

	utils.LogSuccess("programs", "get_all", "Successfully fetched all programs", requestID, 0, map[string]any{
		"count": len(programs),
	})
	return programs, nil
}

func (s *programService) GetByID(id uint, requestID string) (*entities.Program, error) {
	utils.LogInfo("programs", "get_by_id", "Fetching program by ID", requestID, 0, map[string]any{
		"program_id": id,
	})

	program, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("programs", "get_by_id", "Program not found", requestID, 0, map[string]any{
				"program_id": id,
			})
			return nil, errors.New("program not found")
		}
		utils.LogError("programs", "get_by_id", "Failed to fetch program: "+err.Error(), requestID, 0, map[string]any{
			"program_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("programs", "get_by_id", "Successfully fetched program", requestID, 0, map[string]any{
		"program_id":   program.ID,
		"program_name": program.ProgramName,
	})
	return &program, nil
}

func (s *programService) GetByName(name string, requestID string) (*entities.Program, error) {
	utils.LogInfo("programs", "get_by_name", "Fetching program by name", requestID, 0, map[string]any{
		"program_name": name,
	})

	program, err := s.repo.FindByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("programs", "get_by_name", "Program not found", requestID, 0, map[string]any{
				"program_name": name,
			})
			return nil, errors.New("program not found")
		}
		utils.LogError("programs", "get_by_name", "Failed to fetch program: "+err.Error(), requestID, 0, map[string]any{
			"program_name": name,
		})
		return nil, err
	}

	utils.LogSuccess("programs", "get_by_name", "Successfully fetched program", requestID, 0, map[string]any{
		"program_id":   program.ID,
		"program_name": program.ProgramName,
	})
	return &program, nil
}

func (s *programService) Create(input CreateProgramInput, requestID string, userID uint) (*entities.Program, error) {
	utils.LogInfo("programs", "create", "Attempting to create new program", requestID, userID, map[string]any{
		"program_name": input.ProgramName,
	})

	// Check if program with same name already exists
	_, err := s.repo.FindByName(input.ProgramName)
	if err == nil {
		utils.LogWarning("programs", "create", "Program with this name already exists", requestID, userID, map[string]any{
			"program_name": input.ProgramName,
		})
		return nil, errors.New("program with this name already exists")
	}

	program := &entities.Program{
		ProgramName:  input.ProgramName,
		Description:  input.Description,
		ImageProgram: input.ImageProgram,
	}

	if err := s.repo.Create(program); err != nil {
		utils.LogError("programs", "create", "Failed to create program: "+err.Error(), requestID, userID, map[string]any{
			"program_name": input.ProgramName,
		})
		return nil, err
	}

	utils.LogSuccess("programs", "create", "Program created successfully", requestID, userID, map[string]any{
		"program_id":   program.ID,
		"program_name": program.ProgramName,
	})
	return program, nil
}

func (s *programService) Update(id uint, input UpdateProgramInput, requestID string, userID uint) (*entities.Program, error) {
	utils.LogInfo("programs", "update", "Attempting to update program", requestID, userID, map[string]any{
		"program_id": id,
	})

	program, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("programs", "update", "Program not found", requestID, userID, map[string]any{
				"program_id": id,
			})
			return nil, errors.New("program not found")
		}
		utils.LogError("programs", "update", "Failed to fetch program: "+err.Error(), requestID, userID, map[string]any{
			"program_id": id,
		})
		return nil, err
	}

	// Update only provided fields
	if input.ProgramName != nil {
		// Check if new name already exists (if different from current)
		if *input.ProgramName != program.ProgramName {
			existing, _ := s.repo.FindByName(*input.ProgramName)
			if existing.ID != 0 {
				utils.LogWarning("programs", "update", "Program with this name already exists", requestID, userID, map[string]any{
					"program_name": *input.ProgramName,
				})
				return nil, errors.New("program with this name already exists")
			}
		}
		program.ProgramName = *input.ProgramName
	}
	if input.Description != nil {
		program.Description = *input.Description
	}
	if input.ImageProgram != nil {
		program.ImageProgram = *input.ImageProgram
	}

	if err := s.repo.Update(program); err != nil {
		utils.LogError("programs", "update", "Failed to update program: "+err.Error(), requestID, userID, map[string]any{
			"program_id": id,
		})
		return nil, err
	}

	utils.LogSuccess("programs", "update", "Program updated successfully", requestID, userID, map[string]any{
		"program_id":   program.ID,
		"program_name": program.ProgramName,
	})
	return &program, nil
}

func (s *programService) Delete(id uint, requestID string, userID uint) error {
	utils.LogInfo("programs", "delete", "Attempting to delete program", requestID, userID, map[string]any{
		"program_id": id,
	})

	program, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarning("programs", "delete", "Program not found", requestID, userID, map[string]any{
				"program_id": id,
			})
			return errors.New("program not found")
		}
		utils.LogError("programs", "delete", "Failed to fetch program: "+err.Error(), requestID, userID, map[string]any{
			"program_id": id,
		})
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		utils.LogError("programs", "delete", "Failed to delete program: "+err.Error(), requestID, userID, map[string]any{
			"program_id": id,
		})
		return err
	}

	utils.LogSuccess("programs", "delete", "Program deleted successfully", requestID, userID, map[string]any{
		"program_id":   id,
		"program_name": program.ProgramName,
	})
	return nil
}
