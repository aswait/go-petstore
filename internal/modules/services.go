package modules

import (
	"github.com/go-chi/jwtauth"
	pet "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/pet/service"
	store "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/store/service"
	user "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/user/service"
)

type Services struct {
	User  user.Userer
	Store store.Storer
	Pet   pet.Peter
}

func NewServices(storages Storages, tokenAuth *jwtauth.JWTAuth) *Services {
	return &Services{
		User:  user.NewUserService(storages.User, tokenAuth),
		Store: store.NewStoreService(storages.Store),
		Pet:   pet.NewPetService(storages.Pet),
	}
}
