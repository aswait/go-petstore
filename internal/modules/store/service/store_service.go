package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/store/repository"
)

type Storer interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id string) (models.OrderResponse, error)
	DeleteOrder(ctx context.Context, id string) error
	Inventory(ctx context.Context) (models.PetsStatuses, error)

	StatusCheck(status string) error

	Decode(r io.ReadCloser, data interface{}) error
	URLParam(r *http.Request, param string) string
}

type StoreService struct {
	storage repository.StoreRepository
}

func NewStoreService(storage repository.StoreRepository) *StoreService {
	return &StoreService{
		storage: storage,
	}
}

func (s *StoreService) Decode(r io.ReadCloser, data interface{}) error {
	return json.NewDecoder(r).Decode(data)
}

func (s StoreService) StatusCheck(status string) error {
	statuses := []string{"placed", "approved", "delivered"}
	for _, stat := range statuses {
		if stat == status {
			return nil
		}
	}
	return fmt.Errorf("invalid status: %s", status)
}

func (s *StoreService) CreateOrder(ctx context.Context, order models.Order) error {
	return s.storage.CreateOrder(ctx, order)
}

func (s *StoreService) URLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func (s *StoreService) GetByID(ctx context.Context, id string) (models.OrderResponse, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return models.OrderResponse{}, err
	}

	dbOrder, err := s.storage.GetByID(ctx, intId)

	order := models.OrderResponse{
		ID:       dbOrder.ID,
		PetID:    dbOrder.PetID,
		Quantity: dbOrder.Quantity,
		ShipDate: dbOrder.ShipDate,
		Status:   dbOrder.Status,
		Complete: dbOrder.Complete,
	}

	return order, err
}

func (s *StoreService) DeleteOrder(ctx context.Context, id string) error {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	err = s.storage.DeleteOrder(ctx, intId)

	return err
}

func (s *StoreService) Inventory(ctx context.Context) (models.PetsStatuses, error) {
	return s.storage.Inventory(ctx)
}
