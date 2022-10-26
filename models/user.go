package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      string `gorm:"uniqueIndex"`
	UserName    string
	AccessToken string
	Email       string
}

func GetUserByUserId(db *gorm.DB, userId string) (*User, error) {
	var user User
	result := db.Where("user_id = ?", userId).First(&user)
	err := result.Error
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.Where("id = ?", id).First(&user)
	err := result.Error
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveOrUpdateByUserId(db *gorm.DB, user *User) error {
	oldUser, err := GetUserByUserId(db, user.UserID)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if oldUser == nil {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	} else {
		user.ID = oldUser.ID
		if err := db.Save(&user).Error; err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}
	}
	return nil
}
