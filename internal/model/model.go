package model

import (
	"encoding/json"
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

// Task Models
type TaskStatus string

const (
	StatusPending   TaskStatus = "PENDING"
	StatusRunning   TaskStatus = "RUNNING"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusFailed    TaskStatus = "FAILED"
)

type TaskType string

const (
	AIText     TaskType = "AI_TEXT_GEN"
	AIImageGen TaskType = "AI_IMAGE_GEN"
)

type Task struct {
	TaskId    int             `json:"taskId"`
	UserId    int             `json:"UserId"`
	TaskType  TaskType        `json:"taskType,omitempty"`
	Status    TaskStatus      `json:"status,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTask struct {
	TaskType TaskType        `json:"taskType"`
	Payload  json.RawMessage `json:"payload"`
}
