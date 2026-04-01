package model

import "time"

type User struct {
	UserId    int
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt time.Time
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type JWTModel struct {
	UserId *int
	Email  *string
}
