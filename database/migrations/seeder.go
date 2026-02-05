package migrations

import (
	"log"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/gorm"
)

// SeedSubtests seeds the 7 fixed subtests for Try Out.
// This should be called once after migration.
// Order of insertion determines the subtest order during Try Out.
func SeedSubtests(db *gorm.DB) error {
	subtests := []entities.Subtest{
		{
			Code:             "PU",
			Name:             "Penalaran Umum",
			QuestionCount:    30,
			TimeLimitSeconds: 1800, // 30 minutes
			MaxScore:         886.12,
		},
		{
			Code:             "PBM",
			Name:             "Penalaran Matematika",
			QuestionCount:    20,
			TimeLimitSeconds: 1500, // 25 minutes
			MaxScore:         857.13,
		},
		{
			Code:             "PPU",
			Name:             "Pemahaman Bacaan dan Menulis",
			QuestionCount:    20,
			TimeLimitSeconds: 900, // 15 minutes
			MaxScore:         875.96,
		},
		{
			Code:             "PK",
			Name:             "Pengetahuan Kuantitatif",
			QuestionCount:    20,
			TimeLimitSeconds: 1200, // 20 minutes
			MaxScore:         1000.00,
		},
		{
			Code:             "LBI",
			Name:             "Literasi Bahasa Indonesia",
			QuestionCount:    30,
			TimeLimitSeconds: 2550, // 42.5 minutes
			MaxScore:         885.75,
		},
		{
			Code:             "LBE",
			Name:             "Literasi Bahasa Inggris",
			QuestionCount:    20,
			TimeLimitSeconds: 1200, // 20 minutes
			MaxScore:         870.76,
		},
		{
			Code:             "PM",
			Name:             "Penalaran Matematika",
			QuestionCount:    20,
			TimeLimitSeconds: 2550, // 42.5 minutes
			MaxScore:         1000.00,
		},
	}

	for _, subtest := range subtests {
		// Use FirstOrCreate to avoid duplicates
		result := db.Where("code = ?", subtest.Code).FirstOrCreate(&subtest)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			log.Printf("✅ Seeded subtest: %s (%s)", subtest.Name, subtest.Code)
		} else {
			log.Printf("⏭️  Subtest already exists: %s (%s)", subtest.Name, subtest.Code)
		}
	}

	return nil
}
