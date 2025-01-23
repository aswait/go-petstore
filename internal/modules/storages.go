package modules

import (
	"gorm.io/gorm"
	pet "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/pet/repository"
	store "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/store/repository"
	user "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/user/repository"
)

type Storages struct {
	User  user.UserRepository
	Store store.StoreRepository
	Pet   pet.PetRepository
}

func NewStorages(adapter *gorm.DB) *Storages {
	return &Storages{
		User:  user.NewUserStorage(adapter),
		Store: store.NewStoreStorage(adapter),
		Pet:   pet.NewPetStorage(adapter),
	}
}
