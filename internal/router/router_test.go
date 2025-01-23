package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules"
)

type MockUserController struct {
}

func (m *MockUserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) CreateWithListAndArray(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockUserController) GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type MockPetController struct {
}

func (m *MockPetController) CreatePet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) GetByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) FindByStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) FindByTags(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) UpdateByPetId(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) DeleteByPetId(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) UpdatePet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockPetController) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

type MockStoreController struct {
}

func (m *MockStoreController) Order(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockStoreController) Inventory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}
func (m *MockStoreController) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func (m *MockStoreController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestNewRouter(t *testing.T) {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	mockUserController := MockUserController{}
	mockStoreController := MockStoreController{}
	mockPetController := MockPetController{}

	controllers := &modules.Controllers{
		User:  &mockUserController,
		Store: &mockStoreController,
		Pet:   &mockPetController,
	}

	router := NewRouter(controllers, tokenAuth)

	tests := []struct {
		method string
		route  string
	}{
		{"POST", "/user"},
		{"POST", "/user/createWithArray"},
		{"POST", "/user/createWithList"},
		{"GET", "/user/login"},
		{"GET", "/user/logout"},
		{"GET", "/user/testuser"},
		{"PUT", "/user/testuser"},
		{"DELETE", "/user/testuser"},
		{"POST", "/store/order"},
		{"GET", "/store/order/1"},
		{"DELETE", "/store/order/1"},
		{"GET", "/store/inventory"},
		{"POST", "/pet"},
		{"PUT", "/pet"},
		{"GET", "/pet/1"},
		{"POST", "/pet/1"},
		{"DELETE", "/pet/1"},
		{"POST", "/pet/1/uploadImage"},
		{"GET", "/pet/findByStatus"},
		{"GET", "/pet/findByTags"},
		{"GET", "/swagger/"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.route, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if strings.Split(test.route, "/")[1] == "pet" || (test.route == "/user/testuser" && test.method != "GET") || test.route == "/store/inventory" {
			status := rr.Code
			if status != http.StatusForbidden {
				t.Errorf("handler for %s %s returned wrong status code: got %v want %v",
					test.method, test.route, status, http.StatusForbidden)
			}
		} else if test.route == "/swagger/" {
			status := rr.Code
			if status != http.StatusNotFound {
				t.Errorf("handler for %s %s returned wrong status code: got %v want %v",
					test.method, test.route, status, http.StatusNotFound)
			}
		} else if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler for %s %s returned wrong status code: got %v want %v",
				test.method, test.route, status, http.StatusOK)
		}
	}
}
