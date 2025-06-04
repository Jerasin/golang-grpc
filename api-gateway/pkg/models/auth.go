package models

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
