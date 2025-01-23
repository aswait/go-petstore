package controller

import (
	"net/http"

	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/store/service"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/responder"
)

type Data struct {
	Message string `json:"message"`
}

type OrderResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Storer interface {
	Order(w http.ResponseWriter, r *http.Request)
	Inventory(w http.ResponseWriter, r *http.Request)
	GetOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
}

type Store struct {
	service service.Storer
	responder.Responder
}

func NewStore(service service.Storer, responder responder.Responder) *Store {
	return &Store{
		service:   service,
		Responder: responder,
	}
}

func (s *Store) Order(w http.ResponseWriter, r *http.Request) {
	var req models.Order

	err := s.service.Decode(r.Body, &req)
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	err = s.service.StatusCheck(req.Status)
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	err = s.service.CreateOrder(r.Context(), req)
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	s.OutputJSON(w, OrderResponse{
		Success: true,
		Data: Data{
			Message: "Order created successfully",
		},
	})
}

func (s *Store) Inventory(w http.ResponseWriter, r *http.Request) {
	inventory, err := s.service.Inventory(r.Context())
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	s.OutputJSON(w, inventory)
}

func (s *Store) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := s.service.URLParam(r, "orderId")

	order, err := s.service.GetByID(r.Context(), orderID)
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	s.OutputJSON(w, order)
}

func (s *Store) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := s.service.URLParam(r, "orderId")

	err := s.service.DeleteOrder(r.Context(), orderID)
	if err != nil {
		s.Responder.ErrorBadRequest(w, err)
		return
	}

	s.OutputJSON(w, OrderResponse{
		Success: true,
		Data: Data{
			Message: "Order deleted successfully",
		},
	})
}
