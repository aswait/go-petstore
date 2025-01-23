package repository

import (
	"context"

	"gorm.io/gorm"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User, updatedUser models.User) error
	DeleteUser(ctx context.Context, user models.User) error
}

type UserStorage struct {
	adapter *gorm.DB
}

func NewUserStorage(adapter *gorm.DB) *UserStorage {
	return &UserStorage{
		adapter: adapter,
	}
}

func (s *UserStorage) Create(ctx context.Context, user models.User) error {
	result := s.adapter.WithContext(ctx).Create(&user)

	return result.Error
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var users []models.User

	result := s.adapter.WithContext(ctx).Where(&models.User{
		Username:   username,
		UserStatus: 0,
	}).Find(&users)

	if result.Error != nil {
		return models.User{}, result.Error
	}

	for _, user := range users {
		if user.Username == username && user.UserStatus == 0 {
			return user, nil
		}
	}

	return models.User{}, nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, user models.User, updatedUser models.User) error {
	result := s.adapter.WithContext(ctx).Model(&user).Updates(updatedUser)

	return result.Error
}

func (s *UserStorage) DeleteUser(ctx context.Context, user models.User) error {
	result := s.adapter.WithContext(ctx).Model(&user).Update("user_status", 1)

	return result.Error
}
