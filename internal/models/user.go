package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID           string  `gorm:"type:uuid;primaryKey"`
	TelegramUserID int64   `gorm:"type:bigint;unique"`
	DisplayName    string  `gorm:"type:varchar(255)"`
	Balance        float64 `gorm:"type:decimal(10,2)"`
	State          State   `gorm:"type:json"`
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (*UserRepo, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate User: %w", err)
	}

	return &UserRepo{db: db}, nil
}

func (r *UserRepo) GetOrCreateUserByTelegramUserID(telegramUserID int64) (*User, error) {
	user := new(User)

	err := r.db.Where("telegram_user_id = ?", telegramUserID).First(user).Error
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user = &User{
		UUID:           uuid.New().String(),
		TelegramUserID: telegramUserID,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Save(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}
