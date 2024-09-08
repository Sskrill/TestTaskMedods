package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}
type SignUpInput struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

type User struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

func (i SignUpInput) IsValid() error {
	return validate.Struct(i)
}
func (i SignInInput) IsValid() error {
	return validate.Struct(i)
}
