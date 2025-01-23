package repository

import (
	"context"

	"gorm.io/gorm"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
)

type PetRepository interface {
	CreatePet(ctx context.Context, pet models.Pet) error

	UpdatePet(ctx context.Context, pet models.Pet, name string, status string) error
	UpdatePetByModel(ctx context.Context, pet models.Pet, updatedPet models.Pet) error

	GetByName(ctx context.Context, name string) error
	GetByID(ctx context.Context, id int) (models.Pet, error)
	GetCategoryByName(ctx context.Context, category models.Category) (models.Category, error)
	GetTagByName(ctx context.Context, tag models.Tag) (models.Tag, error)
	GetByStatus(ctx context.Context, status string) ([]models.Pet, error)
	GetByTags(ctx context.Context, tags []string) ([]models.Pet, error)

	DeletePet(ctx context.Context, pet models.Pet) error
}

type PetStorage struct {
	adapter *gorm.DB
}

func NewPetStorage(adapter *gorm.DB) *PetStorage {
	return &PetStorage{
		adapter: adapter,
	}
}

func (s *PetStorage) CreatePet(ctx context.Context, pet models.Pet) error {
	result := s.adapter.WithContext(ctx).Create(&pet)
	return result.Error
}

func (s *PetStorage) GetByName(ctx context.Context, name string) error {
	var existingPet models.Pet

	err := s.adapter.WithContext(ctx).Where(&models.Pet{
		Name: name,
	}).First(&existingPet).Error

	return err
}

func (s *PetStorage) GetByID(ctx context.Context, id int) (models.Pet, error) {
	var existingPet models.Pet

	err := s.adapter.WithContext(ctx).Preload("Category").Preload("Tags").Preload("PhotoUrls").Where(&models.Pet{
		ID: id,
	}).First(&existingPet).Error

	return existingPet, err
}

func (s *PetStorage) GetCategoryByName(ctx context.Context, category models.Category) (models.Category, error) {
	var existingCategory models.Category

	err := s.adapter.WithContext(ctx).Where(&models.Category{
		Name: category.Name,
	}).First(&existingCategory).Error

	return existingCategory, err
}

func (s *PetStorage) GetTagByName(ctx context.Context, tag models.Tag) (models.Tag, error) {
	var existingTag models.Tag

	err := s.adapter.WithContext(ctx).Where(&models.Tag{
		Name: tag.Name,
	}).First(&existingTag).Error

	return existingTag, err
}

func (s *PetStorage) GetByStatus(ctx context.Context, status string) ([]models.Pet, error) {
	var existingPets []models.Pet

	err := s.adapter.WithContext(ctx).Preload("Category").Preload("Tags").Where(&models.Pet{
		Status: status,
	}).Find(&existingPets).Error

	return existingPets, err
}

func (s *PetStorage) GetByTags(ctx context.Context, tags []string) ([]models.Pet, error) {
	var existingPets []models.Pet

	err := s.adapter.WithContext(ctx).
		Preload("Category").
		Preload("Tags").
		Joins("JOIN pet_tags ON pet_tags.pet_id = pets.id").
		Joins("JOIN tags ON tags.id = pet_tags.tag_id").
		Where("tags.name IN ?", tags).
		Group("pets.id").
		Having("COUNT(DISTINCT tags.name) = ?", len(tags)).
		Find(&existingPets).Error

	return existingPets, err
}

func (s *PetStorage) UpdatePet(ctx context.Context, pet models.Pet, name string, status string) error {
	err := s.adapter.WithContext(ctx).
		Model(&pet).
		Update("name", name).
		Update("status", status).Error

	return err
}

func (s *PetStorage) DeletePet(ctx context.Context, pet models.Pet) error {
	err := s.adapter.WithContext(ctx).
		Model(&pet).
		Association("Tags").
		Clear()

	if err != nil {
		return err
	}

	err = s.adapter.WithContext(ctx).
		Where("pet_refer_id = ?", pet.ID).
		Delete(&models.PhotoUrl{}).
		Error
	if err != nil {
		return err
	}

	err = s.adapter.WithContext(ctx).Delete(&pet).Error
	return err
}

func (s *PetStorage) UpdatePetByModel(ctx context.Context, pet models.Pet, updatedPet models.Pet) error {
	tx := s.adapter.WithContext(ctx).Begin()

	if len(updatedPet.PhotoUrls) > 0 {
		err := tx.Where("pet_refer_id = ?", pet.ID).Delete(&models.PhotoUrl{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, photo := range updatedPet.PhotoUrls {
			photo.PetReferID = pet.ID
			err := tx.Create(&photo).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(updatedPet.Tags) > 0 {
		err := tx.Model(&pet).Association("Tags").Clear()
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Model(&pet).Association("Tags").Replace(updatedPet.Tags)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err := tx.Model(&pet).Updates(updatedPet).Error
	if err != nil {
		return err
	}

	return tx.Commit().Error
}
