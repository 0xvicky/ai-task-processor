package user

import (
	"ai-task-processor/internal/model"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userDetail model.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, userId int, userUpdateDetail model.UserUpdate) (model.User, error)
	DeleteUser(ctx context.Context, userId int) (model.User, error)
	GetUserById(ctx context.Context, userId int) (model.User, error)
	FetchAllUsers(ctx context.Context) ([]model.User, error)
}
