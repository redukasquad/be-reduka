package migrations

import (
	"fmt"
	"log"
	"os"

	"github.com/redukasquad/be-reduka/database/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
		&entities.ClassSubject{},
		&entities.ClassLesson{},

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
