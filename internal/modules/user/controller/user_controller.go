package controller

import (
	"net/http"
	"time"

	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/user/service"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/responder"
)

type Data struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Userer interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	CreateWithListAndArray(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type User struct {
	service service.Userer
	responder.Responder
}

func NewUser(service service.Userer, responder responder.Responder) *User {
	return &User{
		service:   service,
		Responder: responder,
	}
}

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := u.service.Decode(r.Body, &req)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	// Проверка существования пользователя с таким же юзернеймом
	err = u.service.UserValidation(r.Context(), req.Username)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидация почты
	err = u.service.EmailValidation(req.Email)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидация пароля
	err = u.service.PasswordValidation(req.Password)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидаци номера телефона
	err = u.service.PhoneValidation(req.Phone)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Хэширование пароля
	hpass, err := u.service.PasswordEncryption(req.Password)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}
	req.Password = hpass

	// Создание пользователя в дб
	err = u.service.UserCreate(r.Context(), req)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	u.Responder.OutputJSON(w, UserResponse{
		Success: true,
		Data: Data{
			Message: "user created successfully",
		},
	})
}

func (u *User) CreateWithListAndArray(w http.ResponseWriter, r *http.Request) {
	var req []models.User
	err := u.service.Decode(r.Body, &req)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	var responses []UserResponse

	// Перенаправление запроса для каждого юзера
	for _, user := range req {
		resp, err := u.service.UserRequestRedirection(user)
		if err != nil {
			u.Responder.ErrorBadRequest(w, err)
			return
		}

		var userResp UserResponse

		err = u.service.Decode(resp.Body, &userResp)
		if err != nil {
			u.Responder.ErrorBadRequest(w, err)
			return
		}

		responses = append(responses, userResp)
	}

	u.Responder.OutputJSON(w, responses)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var query models.LoginForm
	err := u.service.DecodeURl(&query, r.URL.Query())
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	user, err := u.service.UserExistenceCheck(r.Context(), query.Username)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	err = u.service.PasswordCheck(r.Context(), query, user)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	token, err := u.service.MakeToken(query.Username)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}
	u.service.SetCookie(w, true, token)
	w.Header().Add("X-Expires-After", time.Now().Add(7*24*time.Hour).String())
	w.Header().Add("X-Rate-Limit", "50")
	w.Write([]byte(token))
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	u.service.SetCookie(w, false, "")

	u.Responder.OutputJSON(w, UserResponse{
		Success: true,
		Data: Data{
			Message: "successful logout",
		},
	})
}

func (u *User) GetUser(w http.ResponseWriter, r *http.Request) {
	username := u.service.URLParam(r, "username")

	user, err := u.service.UserExistenceCheck(r.Context(), username)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	u.Responder.OutputJSON(w, user)
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := u.service.URLParam(r, "username")

	var req models.User
	err := u.service.Decode(r.Body, &req)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	user, err := u.service.UserExistenceCheck(r.Context(), username)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	err = u.service.UserValidation(r.Context(), req.Username)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидация почты
	err = u.service.EmailValidation(req.Email)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидация пароля
	err = u.service.PasswordValidation(req.Password)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Валидаци номера телефона
	err = u.service.PhoneValidation(req.Phone)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}

	// Хэширование пароля
	hpass, err := u.service.PasswordEncryption(req.Password)
	if err != nil {
		u.Responder.OutputJSON(w, UserResponse{
			Success: false,
			Data: Data{
				Message: err.Error(),
			},
		})
		return
	}
	req.Password = hpass

	err = u.service.UpdateUser(r.Context(), user, req)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	u.Responder.OutputJSON(w, UserResponse{
		Success: true,
		Data: Data{
			Message: "user successfully updated",
		},
	})
}

func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := u.service.URLParam(r, "username")

	user, err := u.service.UserExistenceCheck(r.Context(), username)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	err = u.service.DeleteUser(r.Context(), user)
	if err != nil {
		u.Responder.ErrorBadRequest(w, err)
		return
	}

	u.Responder.OutputJSON(w, UserResponse{
		Success: true,
		Data: Data{
			Message: "user successfully deleted",
		},
	})
}
