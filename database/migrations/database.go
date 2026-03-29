package migrations

import (
	"fmt"
	"log"
	"os"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
  "host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
  os.Getenv("DB_HOST"),
  os.Getenv("DB_USER"),
  os.Getenv("DB_PASS"),
  os.Getenv("DB_NAME"),
  os.Getenv("DB_PORT"),
)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	log.Println("Running auto-migrations...")

	log.Println("🚀 Running auto-migrations...")
	err = DB.AutoMigrate(
		// ===== USER & AUTH =====
		&entities.User{},

		// ===== PROGRAM & COURSE =====
		&entities.Program{},
		&entities.Course{},
		&entities.CourseRegistration{},
		&entities.RegistrationQuestion{},
		&entities.RegistrationAnswer{},

		// ===== CLASS & LESSON =====
		&entities.Class{},
		&entities.Lesson{},
		&entities.LessonResource{},

		// ===== TRYOUT =====
		&entities.Subtest{},
		&entities.TryOut{},
		&entities.TutorPermission{},
		&entities.TryOutQuestion{},
		&entities.TryOutRegistration{},
		&entities.TryOutAttempt{},
		&entities.SubtestResult{},
		&entities.UserTryOutAnswer{},

		// ===== UNIVERSITY & TARGET =====
		&entities.University{},
		&entities.UniversityMajor{},
		&entities.UserTarget{},

		// ===== UPLOADS =====
		&entities.Image{},
	)

	if err != nil {
		log.Fatalf("Failed to run auto-migrations: %v", err)
	}
	log.Println("Auto-migrations completed successfully!")

	// Seed master data
	log.Println("🌱 Running seeders...")
	if err := SeedSubtests(DB); err != nil {
		log.Fatalf("Failed to seed subtests: %v", err)
	}
	log.Println("Seeders completed successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
