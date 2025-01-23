package controller

import (
	"net/http"

	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/pet/service"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/responder"
)

type Data struct {
	Message string `json:"message"`
}

type PetResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Peter interface {
	CreatePet(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	FindByStatus(w http.ResponseWriter, r *http.Request)
	FindByTags(w http.ResponseWriter, r *http.Request)
	UpdateByPetId(w http.ResponseWriter, r *http.Request)
	DeleteByPetId(w http.ResponseWriter, r *http.Request)
	UpdatePet(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
}

type Pet struct {
	service service.Peter
	responder.Responder
}

func NewPet(service service.Peter, responder responder.Responder) *Pet {
	return &Pet{
		service:   service,
		Responder: responder,
	}
}

func (p *Pet) CreatePet(w http.ResponseWriter, r *http.Request) {
	var pet models.PetJSON
	err := p.service.Decode(r.Body, &pet)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.StatusCheck(pet.Status)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	dbPet := p.service.PetToDB(pet)

	err = p.service.ExistingPet(r.Context(), dbPet.Name)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.service.ExistingCategory(r.Context(), &dbPet)
	p.service.ExistingTag(r.Context(), &dbPet)

	err = p.service.CreatePet(r.Context(), dbPet)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Success: true,
		Data: Data{
			Message: "pet created successfully",
		},
	})
}

func (p *Pet) GetByID(w http.ResponseWriter, r *http.Request) {
	id := p.service.URLParam(r, "petId")

	pet, err := p.service.GetPetByID(r.Context(), id)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, pet)
}

func (p *Pet) FindByStatus(w http.ResponseWriter, r *http.Request) {
	var query models.StatusForm

	err := p.service.DecodeURl(&query, r.URL.Query())
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	pets, _ := p.service.FindByStatus(r.Context(), query.Statuses)

	p.OutputJSON(w, pets)
}

func (p *Pet) FindByTags(w http.ResponseWriter, r *http.Request) {
	var query models.TagsForm

	err := p.service.DecodeURl(&query, r.URL.Query())
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	pets, err := p.service.FindByTags(r.Context(), query.Tags)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, pets)
}

func (p *Pet) UpdateByPetId(w http.ResponseWriter, r *http.Request) {
	id := p.service.URLParam(r, "petId")

	pet, err := p.service.GetPetByID(r.Context(), id)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	form, err := p.service.ValuesFromForm(r)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.ExistingPet(r.Context(), form.Name)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.StatusCheck(form.Status)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.UpdatePet(r.Context(), pet, form)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Success: true,
		Data: Data{
			Message: "pet updated successfully",
		},
	})
}

func (p *Pet) DeleteByPetId(w http.ResponseWriter, r *http.Request) {
	id := p.service.URLParam(r, "petId")

	pet, err := p.service.GetPetByID(r.Context(), id)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.DeletePet(r.Context(), pet)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Success: true,
		Data: Data{
			Message: "pet deleted successfully",
		},
	})
}

func (p *Pet) UpdatePet(w http.ResponseWriter, r *http.Request) {
	var pet models.PetJSON

	err := p.service.Decode(r.Body, &pet)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	dbPet, err := p.service.GetPetByID(r.Context(), p.service.Itoa(pet.ID))
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	updatedPet := p.service.PetToDB(pet)

	err = p.service.ExistingPet(r.Context(), updatedPet.Name)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.StatusCheck(updatedPet.Status)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	// p.service.ExistingCategory(r.Context(), &updatedPet)
	p.service.ExistingTag(r.Context(), &updatedPet)

	err = p.service.UpdatePetByModel(r.Context(), dbPet, updatedPet)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Success: true,
		Data: Data{
			Message: "pet updated successfully",
		},
	})
}

func (p *Pet) UploadImage(w http.ResponseWriter, r *http.Request) {
	id := p.service.URLParam(r, "petId")

	fileName, err := p.service.FileFromForm(r)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	dbPet, err := p.service.GetPetByID(r.Context(), id)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	err = p.service.AddPetPhotoUrls(r.Context(), dbPet, fileName)
	if err != nil {
		p.Responder.ErrorBadRequest(w, err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Success: true,
		Data: Data{
			Message: "image uploaded successfully",
		},
	})
}
