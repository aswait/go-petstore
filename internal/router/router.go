package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/middleware"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules"
)

func NewRouter(controllers *modules.Controllers, tokenAuth *jwtauth.JWTAuth) http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Route("/user", func(r chi.Router) {
		r.Post("/", controllers.User.CreateUser)
		r.Post("/createWithArray", controllers.User.CreateWithListAndArray)
		r.Post("/createWithList", controllers.User.CreateWithListAndArray)

		r.With().Get("/login", controllers.User.Login)
		r.Get("/logout", controllers.User.Logout)

		r.Get("/{username}", controllers.User.GetUser)

		r.Group(func(r chi.Router) {
			r.Use(middleware.UserUnloggedIn)

			r.Put("/{username}", controllers.User.UpdateUser)
			r.Delete("/{username}", controllers.User.DeleteUser)
		})
	})

	r.Route("/store", func(r chi.Router) {
		r.Route("/order", func(r chi.Router) {
			r.Post("/", controllers.Store.Order)

			r.Route("/{orderId}", func(r chi.Router) {
				r.Get("/", controllers.Store.GetOrder)
				r.Delete("/", controllers.Store.DeleteOrder)
			})
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.UnloggedIn)
			r.Get("/inventory", controllers.Store.Inventory)
		})
	})

	r.Route("/pet", func(r chi.Router) {
		r.Use(middleware.UnloggedIn)
		r.Route("/{petId}", func(r chi.Router) {
			r.Route("/", func(r chi.Router) {
				r.Get("/", controllers.Pet.GetByID)
				r.Post("/", controllers.Pet.UpdateByPetId)

				r.Group(func(r chi.Router) {
					r.Use(middleware.UnloggedInDelete)

					r.Delete("/", controllers.Pet.DeleteByPetId)
				})
			})
			r.Post("/uploadImage", controllers.Pet.UploadImage)
		})
		r.Route("/", func(r chi.Router) {
			r.Post("/", controllers.Pet.CreatePet)
			r.Put("/", controllers.Pet.UpdatePet)
		})
		r.Get("/findByStatus", controllers.Pet.FindByStatus)
		r.Get("/findByTags", controllers.Pet.FindByTags)
	})

	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/swagger/", http.FileServer(http.Dir("/public"))).ServeHTTP(w, r)
	})
	return r
}
