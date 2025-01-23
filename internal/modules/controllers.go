package modules

import (
	pet "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/pet/controller"
	store "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/store/controller"
	user "studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/user/controller"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/responder"
)

type Controllers struct {
	User  user.Userer
	Store store.Storer
	Pet   pet.Peter
}

func NewControllers(services *Services, responder responder.Responder) *Controllers {
	return &Controllers{
		User:  user.NewUser(services.User, responder),
		Store: store.NewStore(services.Store, responder),
		Pet:   pet.NewPet(services.Pet, responder),
	}
}
