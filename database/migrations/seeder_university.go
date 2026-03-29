package migrations

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redukasquad/be-reduka/database/entities"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// SeedUniversitiesAndPrograms reads university and program data from an Excel file
// and seeds them into the database using FirstOrCreate (idempotent).
func SeedUniversitiesAndPrograms(db *gorm.DB) error {
	filePath := "passing_grade_ptn_full.xlsx"

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file '%s': %w", filePath, err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows from sheet '%s': %w", sheetName, err)
	}

	if len(rows) < 2 {
		log.Println("⚠️  Excel file has no data rows, skipping university seeder")
		return nil
	}

	// ==========================================
	// STEP 1: Collect unique universities
	// ==========================================
	universityMap := make(map[string]uint) // Map name -> DB ID
	universityCount := 0

	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}
		if len(row) < 3 {
			continue
		}

		uniName := strings.TrimSpace(row[0])
		if uniName == "" {
			continue
		}

		// Skip if already processed
		if _, exists := universityMap[uniName]; exists {
			continue
		}

		uni := entities.University{
			Name: uniName,
			Type: "PTN", // All data in this file are PTN
		}

		result := db.Where("name = ?", uni.Name).FirstOrCreate(&uni)
		if result.Error != nil {
			return fmt.Errorf("failed to seed university '%s': %w", uniName, result.Error)
		}

		universityMap[uniName] = uni.ID

		if result.RowsAffected > 0 {
			universityCount++
			log.Printf("✅ Seeded university: %s", uniName)
		}
	}

	log.Printf("📊 Universities: %d new, %d total in file", universityCount, len(universityMap))

	// ==========================================
	// STEP 2: Seed programs (prodi) with passing grade
	// ==========================================
	programCount := 0
	skippedCount := 0

	for i, row := range rows {
		if i == 0 { // skip header
			continue
		}
		if len(row) < 3 {
			continue
		}

		uniName := strings.TrimSpace(row[0])
		prodiName := strings.TrimSpace(row[1])
		passingGradeStr := strings.TrimSpace(row[2])

		if uniName == "" || prodiName == "" {
			continue
		}

		universityID, exists := universityMap[uniName]
		if !exists {
			log.Printf("⚠️  University not found in map: %s (skipping: %s)", uniName, prodiName)
			skippedCount++
			continue
		}

		passingGrade, err := strconv.ParseFloat(passingGradeStr, 64)
		if err != nil {
			log.Printf("⚠️  Invalid passing grade for %s - %s: '%s' (skipping)", uniName, prodiName, passingGradeStr)
			skippedCount++
			continue
		}

		program := entities.UniversityMajor{
			UniversityID: universityID,
			Name:         prodiName,
			PassingGrade: passingGrade,
		}

		result := db.Where("university_id = ? AND name = ?", universityID, prodiName).FirstOrCreate(&program)
		if result.Error != nil {
			return fmt.Errorf("failed to seed program '%s' for '%s': %w", prodiName, uniName, result.Error)
		}

		if result.RowsAffected > 0 {
			programCount++
		}
	}

	log.Printf("📊 Programs: %d new, %d skipped", programCount, skippedCount)
	log.Printf("🎉 University seeder completed!")

	return nil
}
