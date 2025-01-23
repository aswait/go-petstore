package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
)

type StoreRepository interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int) (models.Order, error)
	DeleteOrder(ctx context.Context, id int) error

	Inventory(ctx context.Context) (models.PetsStatuses, error)
}

type StoreStorage struct {
	adapter *gorm.DB
}

func NewStoreStorage(adapter *gorm.DB) *StoreStorage {
	return &StoreStorage{
		adapter: adapter,
	}
}

func (s *StoreStorage) CreateOrder(ctx context.Context, order models.Order) error {
	var pet models.Pet

	err := s.adapter.WithContext(ctx).First(&pet, order.PetID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pet with ID %d does not exist", order.PetID)
		}
		return err
	}

	order.Pet = pet

	err = s.adapter.WithContext(ctx).Create(&order).Error

	return err
}

func (s *StoreStorage) GetByID(ctx context.Context, id int) (models.Order, error) {
	var existingOrder models.Order

	err := s.adapter.WithContext(ctx).
		Where(&models.Order{
			ID: id,
		}).First(&existingOrder).Error

	return existingOrder, err
}

func (s *StoreStorage) DeleteOrder(ctx context.Context, id int) error {
	order, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.adapter.Delete(&order).Error

	return err
}

func (s *StoreStorage) Inventory(ctx context.Context) (models.PetsStatuses, error) {
	var statuses models.PetsStatuses
	var count int64

	err := s.adapter.WithContext(ctx).
		Model(&models.Pet{}).
		Where("status = ?", "available").
		Count(&count).Error
	if err != nil {
		return models.PetsStatuses{}, err
	}

	statuses.Available = int(count)

	err = s.adapter.WithContext(ctx).
		Model(&models.Pet{}).
		Where("status = ?", "pending").
		Count(&count).Error
	if err != nil {
		return models.PetsStatuses{}, err
	}

	statuses.Pending = int(count)

	err = s.adapter.WithContext(ctx).
		Model(&models.Pet{}).
		Where("status = ?", "sold").
		Count(&count).Error
	if err != nil {
		return models.PetsStatuses{}, err
	}

	statuses.Sold = int(count)

	return statuses, nil
}
