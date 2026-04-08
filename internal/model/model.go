package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserId    int    `json:"userId,omitempty"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Role      string `json:"role,omitempty"`
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

type JwtAuthRes struct {
	UserId   int
	UserRole string
	JwtToken string
}

type JwtClaims struct {
	UserId   int    `json:"userId"`
	UserRole string `json:"userRole"`
	jwt.RegisteredClaims
}

type Task struct {
	TaskId    int    `json:"taskId"`
	UserId    int    `json:"UserId"`
	TaskType  string `json:"taskType,omitempty"`
	Status    string `json:"status,omitempty"`
	Result    any    `json:"result,omitempty"`
	CreatedAt time.Time
}
