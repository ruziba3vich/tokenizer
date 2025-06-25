package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	lgg "github.com/ruziba3vich/prodonik_lgger"
	"github.com/ruziba3vich/tokenizer/internal/models"
	"gorm.io/gorm"
)

type (
	Service struct {
		db     *gorm.DB
		logger *lgg.Logger
	}
)

func NewService(db *gorm.DB, logger *lgg.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (h *Service) GenerateOneTimeLink() (string, error) {
	keyBytes := make([]byte, 8)
	if _, err := rand.Read(keyBytes); err != nil {
		h.logger.Errorf("failed to generate key: %s", err.Error())
		return "", err
	}
	key := hex.EncodeToString(keyBytes)

	link := models.OneTimeLink{Key: key, Used: false}
	if err := h.db.Create(&link).Error; err != nil {
		h.logger.Errorf("failed to store one-time key: %s", err.Error())
		return "", err
	}

	url := fmt.Sprintf("http://localhost:5174/key/%s", key)
	return url, nil
}

func (s *Service) CreateUser(payload *models.RegisterPayload) error {
	var link models.OneTimeLink
	if err := s.db.Where("key = ? AND used = false", payload.Key).First(&link).Error; err != nil {
		return fmt.Errorf("invalid or used key")
	}

	user := models.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Phone:     payload.Phone,
		Username:  payload.Username,
		Password:  payload.Password,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		link.Used = true
		if err := tx.Save(&link).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *Service) GetUserByEmailAndPassword(email, password string) (*models.User, error) {
	var user models.User

	// Only search by email
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Then compare password manually
	if user.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
