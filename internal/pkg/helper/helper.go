package helper

import (
	"fmt"
	"log"
	"time"

	"github.com/ruziba3vich/tokenizer/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dsn = "host=168.119.255.188 user=vinturm_user password=vinturm_pwd dbname=vinturm_db port=5432 sslmode=disable"
)

func GenerateTimeUUID() string {
	now := time.Now()
	timeComponent := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond())
	if len(timeComponent) > 32 {
		timeComponent = timeComponent[:32]
	}
	for len(timeComponent) < 32 {
		timeComponent += "0"
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		timeComponent[0:8],
		timeComponent[8:12],
		timeComponent[12:16],
		timeComponent[16:20],
		timeComponent[20:32])
}

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.OneTimeLink{}, &models.User{})
	if err != nil {
		log.Fatalf("AutoMigration failed: %v", err)
		return nil, err
	}

	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)").Error
	if err != nil {
		log.Printf("Failed to create email index: %v", err)
		return nil, err
	}

	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username)").Error
	if err != nil {
		log.Printf("Failed to create username index: %v", err)
		return nil, err
	}

	log.Println("Migration and indexing successful")
	return db, nil
}
