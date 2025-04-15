package types

import "apps/ecosystem/core/models"

type UserRegisterPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserLoginPayload struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserAuthResponse struct {
	Token string
	User models.UserModel
}
