package database

import (
	"bamboocians/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err = DB.AutoMigrate(
		&models.User{},
		&models.OrganizerProfile{},
		&models.TalentProfile{},
		&models.VenueProfile{},
		&models.Event{},
		&models.Booking{},
		&models.Contract{},
		&models.Guest{},
		&models.Ticket{},
		&models.Message{},
	); err != nil {
		log.Fatalf("auto-migration failed: %v", err)
	}

	log.Println("database connected and migrated")
}
