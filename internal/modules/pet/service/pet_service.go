package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/form"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/pet/repository"
)

type Peter interface {
	ExistingPet(ctx context.Context, name string) error
	ExistingTag(ctx context.Context, pet *models.Pet)
	ExistingCategory(ctx context.Context, pet *models.Pet)
	CreatePet(ctx context.Context, pet models.Pet) error
	GetPetByID(ctx context.Context, id string) (models.Pet, error)
	FindByStatus(ctx context.Context, status []string) ([]models.Pet, error)
	FindByTags(ctx context.Context, tags []string) ([]models.Pet, error)
	UpdatePet(ctx context.Context, pet models.Pet, form models.PetIdForm) error
	DeletePet(ctx context.Context, pet models.Pet) error
	UpdatePetByModel(ctx context.Context, pet models.Pet, updatedPet models.Pet) error
	AddPetPhotoUrls(ctx context.Context, pet models.Pet, url string) error

	StatusCheck(status string) error
	PetToDB(pet models.PetJSON) models.Pet
	Itoa(id int) string

	Decode(r io.ReadCloser, data interface{}) error
	DecodeURl(params interface{}, values url.Values) error
	URLParam(r *http.Request, param string) string
	ValuesFromForm(r *http.Request) (models.PetIdForm, error)
	FileFromForm(r *http.Request) (string, error)
}

type PetService struct {
	storage repository.PetRepository
}

func NewPetService(storage repository.PetRepository) *PetService {
	return &PetService{
		storage: storage,
	}
}

func (s *PetService) Decode(r io.ReadCloser, data interface{}) error {
	return json.NewDecoder(r).Decode(data)
}

func (s *PetService) StatusCheck(status string) error {
	statuses := []string{"available", "pending", "sold"}
	for _, stat := range statuses {
		if stat == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s", status)
}

func (s *PetService) PetToDB(pet models.PetJSON) models.Pet {
	petDB := models.Pet{
		Category: pet.Category,
		Name:     pet.Name,
		Tags:     pet.Tags,
		Status:   pet.Status,
	}
	var urls []models.PhotoUrl
	for _, photourl := range pet.PhotoUrls {
		urls = append(urls, models.PhotoUrl{
			PhotoUrl: photourl,
		})
	}
	petDB.PhotoUrls = urls
	return petDB
}

func (s *PetService) CreatePet(ctx context.Context, pet models.Pet) error {
	return s.storage.CreatePet(ctx, pet)
}

func (s *PetService) ExistingPet(ctx context.Context, name string) error {
	err := s.storage.GetByName(ctx, name)
	if err == nil {
		return fmt.Errorf("a pet with that name already exists")
	}
	return nil
}

func (s *PetService) ExistingCategory(ctx context.Context, pet *models.Pet) {
	category, err := s.storage.GetCategoryByName(ctx, pet.Category)
	if err != nil {
		return
	}
	pet.Category = category
}

func (s *PetService) ExistingTag(ctx context.Context, pet *models.Pet) {
	var tags []models.Tag
	for _, tag := range pet.Tags {
		existTag, err := s.storage.GetTagByName(ctx, tag)
		if err == nil {
			tags = append(tags, existTag)
		} else {
			tags = append(tags, tag)
		}
	}
	pet.Tags = tags
	return
}

func (s *PetService) URLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func (s *PetService) GetPetByID(ctx context.Context, id string) (models.Pet, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return models.Pet{}, err
	}

	pet, err := s.storage.GetByID(ctx, intId)
	if err != nil {
		return models.Pet{}, fmt.Errorf("pet with that id does not exist: %v", id)
	}

	return pet, nil
}

func (s *PetService) DecodeURl(params interface{}, values url.Values) error {
	return form.NewDecoder().Decode(&params, values)
}

func (s *PetService) FindByStatus(ctx context.Context, statuses []string) ([]models.Pet, error) {
	var pets []models.Pet

	for _, status := range statuses {
		dbPets, err := s.storage.GetByStatus(ctx, status)
		if err != nil {
			return pets, err
		}
		pets = append(pets, dbPets...)
	}
	return pets, nil
}

func (s *PetService) FindByTags(ctx context.Context, tags []string) ([]models.Pet, error) {
	pets, err := s.storage.GetByTags(ctx, tags)
	if err != nil {
		return []models.Pet{}, err
	}
	return pets, nil
}

func (s *PetService) ValuesFromForm(r *http.Request) (models.PetIdForm, error) {
	err := r.ParseForm()
	if err != nil {
		return models.PetIdForm{}, err
	}

	petIdForm := models.PetIdForm{
		Name:   r.FormValue("name"),
		Status: r.FormValue("status"),
	}

	return petIdForm, nil
}

func (s *PetService) FileFromForm(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}
	formData := r.MultipartForm.File
	fileName := formData["file"][0].Filename

	return fileName, nil
}

func (s *PetService) UpdatePet(ctx context.Context, pet models.Pet, form models.PetIdForm) error {
	err := s.storage.UpdatePet(ctx, pet, form.Name, form.Status)

	return err
}

func (s *PetService) DeletePet(ctx context.Context, pet models.Pet) error {
	err := s.storage.DeletePet(ctx, pet)

	return err
}

func (s *PetService) Itoa(id int) string {
	return strconv.Itoa(id)
}

func (s *PetService) UpdatePetByModel(ctx context.Context, pet models.Pet, updatedPet models.Pet) error {
	return s.storage.UpdatePetByModel(ctx, pet, updatedPet)
}

func (s *PetService) AddPetPhotoUrls(ctx context.Context, pet models.Pet, url string) error {
	urls := pet.PhotoUrls
	fmt.Println(urls)
	photoUrl := models.PhotoUrl{
		PhotoUrl:   fmt.Sprintf("/pets/%v/%s", pet.ID, url),
		PetReferID: pet.ID,
	}

	urls = append(urls, photoUrl)
	updatedPet := pet
	updatedPet.PhotoUrls = urls

	return s.storage.UpdatePetByModel(ctx, pet, updatedPet)
}
