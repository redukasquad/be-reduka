package migrations

import (
	"log"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

// SeedUniversitiesAndPrograms seeds university and program data.
// TODO: Replace placeholder data with actual data from external source (JSON/CSV).
func SeedUniversitiesAndPrograms(db *gorm.DB) error {
	// ==========================================
	// STEP 1: Seed Universities
	// ==========================================
	universities := []entities.University{
		// TODO: Add universities here
		// Format:
		// {Name: "Universitas Indonesia", Type: "PTN"},
		// {Name: "Institut Teknologi Bandung", Type: "PTN"},
		// {Name: "Binus University", Type: "PTS"},
	}

	universityMap := make(map[string]uint) // Map name -> ID for FK reference

	for _, uni := range universities {
		result := db.Where("name = ?", uni.Name).FirstOrCreate(&uni)
		if result.Error != nil {
			return result.Error
		}
		universityMap[uni.Name] = uni.ID

		if result.RowsAffected > 0 {
			log.Printf("✅ Seeded university: %s (%s)", uni.Name, uni.Type)
		}
	}

	// ==========================================
	// STEP 2: Seed University Programs
	// ==========================================
	type ProgramData struct {
		UniversityName string
		ProgramName    string
		PassingGrade   float64
	}

	programs := []ProgramData{
		// TODO: Add programs here
		// Format:
		// {UniversityName: "Universitas Indonesia", ProgramName: "Teknik Informatika", PassingGrade: 750.50},
		// {UniversityName: "Universitas Indonesia", ProgramName: "Kedokteran", PassingGrade: 850.25},
		// {UniversityName: "Institut Teknologi Bandung", ProgramName: "STEI", PassingGrade: 800.00},
	}

	for _, prog := range programs {
		universityID, exists := universityMap[prog.UniversityName]
		if !exists {
			log.Printf("⚠️  University not found: %s (skipping program: %s)", prog.UniversityName, prog.ProgramName)
			continue
		}

		program := entities.UniversityMajor{
			UniversityID: universityID,
			Name:         prog.ProgramName,
			PassingGrade: prog.PassingGrade,
		}

		result := db.Where("university_id = ? AND name = ?", universityID, prog.ProgramName).FirstOrCreate(&program)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			log.Printf("✅ Seeded program: %s - %s (PG: %.2f)", prog.UniversityName, prog.ProgramName, prog.PassingGrade)
		}
	}

	return nil
}
