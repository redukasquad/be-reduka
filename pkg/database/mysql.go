package database

import (
	"fmt"
	"log"

	"github.com/redukasquad/be-reduka/configs"
	"github.com/redukasquad/be-reduka/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(cfg *configs.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.User{},
	); err != nil {
		return err
	}
	return nil
}
