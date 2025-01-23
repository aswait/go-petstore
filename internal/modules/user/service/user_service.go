package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-passwd/validator"
	"github.com/go-playground/form"
	"golang.org/x/crypto/bcrypt"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/models"
	"studentgit.kata.academy/ponomarenko.100299/go-petstore/internal/modules/user/repository"
)

type Userer interface {
	UserCreate(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User, updatedUser models.User) error
	DeleteUser(ctx context.Context, user models.User) error
	UserExistenceCheck(ctx context.Context, username string) (models.User, error)

	UserValidation(ctx context.Context, username string) error
	EmailValidation(email string) error
	PasswordValidation(password string) error
	PhoneValidation(phone string) error

	PasswordCheck(ctx context.Context, query models.LoginForm, user models.User) error
	PasswordEncryption(password string) (string, error)

	UserRequestRedirection(user models.User) (*http.Response, error)
	DecodeURl(params *models.LoginForm, values url.Values) error
	URLParam(r *http.Request, param string) string
	MakeToken(name string) (string, error)
	SetCookie(w http.ResponseWriter, login bool, value string)

	Decode(r io.ReadCloser, data interface{}) error
}

type UserService struct {
	storage   repository.UserRepository
	tokenAuth *jwtauth.JWTAuth
}

func NewUserService(storage repository.UserRepository, tokenAuth *jwtauth.JWTAuth) *UserService {
	return &UserService{
		storage:   storage,
		tokenAuth: tokenAuth,
	}
}

func (s *UserService) Decode(r io.ReadCloser, data interface{}) error {
	return json.NewDecoder(r).Decode(data)
}

func (s *UserService) UserCreate(ctx context.Context, user models.User) error {
	return s.storage.Create(ctx, user)
}

func (s *UserService) UserValidation(ctx context.Context, username string) error {
	user, err := s.storage.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.ID != 0 {
		return fmt.Errorf("username not available")
	}
	return nil
}

func (s *UserService) EmailValidation(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func (s *UserService) PasswordValidation(password string) error {
	passwordValidator := validator.New(
		validator.MinLength(5, errors.New("password is to short")),
		validator.MaxLength(16, errors.New("password is to long")),
		validator.ContainsOnly(
			"abcdefghijklmnopqrstuvwxyz-.@!$&",
			errors.New("password contains invalid characters"),
		),
	)
	return passwordValidator.Validate(password)
}

func (s *UserService) PhoneValidation(phone string) error {
	pattern := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(pattern)
	phoneNumber := re.Find([]byte(phone))
	if string(phoneNumber) != phone {
		return errors.New("invalid phone number")
	}
	return nil
}

func (s *UserService) PasswordEncryption(password string) (string, error) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hpass), err
}

func (s *UserService) UserRequestRedirection(user models.User) (*http.Response, error) {
	encodedUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/user", bytes.NewBuffer(encodedUser))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (s *UserService) DecodeURl(params *models.LoginForm, values url.Values) error {
	return form.NewDecoder().Decode(&params, values)
}

func (s *UserService) UserExistenceCheck(ctx context.Context, username string) (models.User, error) {
	user, err := s.storage.GetByUsername(ctx, username)
	if err != nil {
		return models.User{}, err
	}

	if user.ID == 0 {
		return models.User{}, fmt.Errorf("user does not exist: %s", username)
	}

	return user, nil
}

func (s *UserService) PasswordCheck(ctx context.Context, query models.LoginForm, user models.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(query.Password))
	if err != nil {
		return fmt.Errorf("wrong password")
	}
	return nil
}

func (s *UserService) MakeToken(name string) (string, error) {
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{"username": name})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserService) SetCookie(w http.ResponseWriter, login bool, value string) {
	if login {
		http.SetCookie(w, &http.Cookie{
			HttpOnly: true,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			SameSite: http.SameSiteLaxMode,
			Name:     "jwt",
			Value:    value,
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    ""},
	)
	return
}

func (s *UserService) URLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User, updatedUser models.User) error {
	return s.storage.UpdateUser(ctx, user, updatedUser)
}

func (s *UserService) DeleteUser(ctx context.Context, user models.User) error {
	return s.storage.DeleteUser(ctx, user)
}
