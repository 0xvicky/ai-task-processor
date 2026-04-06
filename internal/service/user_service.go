package service

import (
	apperrors "ai-task-processor/internal/apperrors"

	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/utils"
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, newUser model.User) (model.JwtAuthRes, error) {
	_, fetchErr := s.repo.GetUserByEmail(ctx, newUser.Email)

	if fetchErr == nil {
		return model.JwtAuthRes{}, apperrors.ErrUserAlreadyExists
	}

	if !errors.Is(fetchErr, sql.ErrNoRows) {
		return model.JwtAuthRes{}, apperrors.ErrInternal
	}

	//get the password out of body
	plainPass := newUser.Password
	//hash the password
	hashPass, hashErr := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)

	if hashErr != nil {
		return model.JwtAuthRes{}, apperrors.ErrInternal
	}
	//save the string(hash) back to struct
	newUser.Password = string(hashPass)
	//passing the user details to the db-repository
	userId, creationErr := s.repo.CreateUser(ctx, newUser)
	if creationErr != nil {
		if errors.Is(creationErr, context.DeadlineExceeded) {
			return model.JwtAuthRes{}, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(creationErr, context.Canceled) {
			return model.JwtAuthRes{}, apperrors.ErrCanceled
		}
		return model.JwtAuthRes{}, apperrors.ErrInternal
	}

	token, jwtErr := utils.JWTInit(userId, newUser.Role)

	if jwtErr != nil {
		return model.JwtAuthRes{}, apperrors.ErrInternal
	}
	newUserRes := model.JwtAuthRes{
		UserId:   userId,
		UserRole: newUser.Role,
		JwtToken: token,
	}

	return newUserRes, nil
}

func (s *UserService) LoginUser(ctx context.Context, userLoginInfo model.UserLogin) (model.JwtAuthRes, error) {
	//check if user exist or not, if exist fetch the user details
	userInfo, fetchErr := s.repo.GetUserByEmail(ctx, userLoginInfo.Email)
	if fetchErr != nil {
		if errors.Is(fetchErr, sql.ErrNoRows) {
			return model.JwtAuthRes{}, apperrors.ErrInvalidCredentials
		}

		return model.JwtAuthRes{}, apperrors.ErrInternal
	}
	// //if exist, then hash the password and compare with stored hash pass

	hashErr := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(userLoginInfo.Password))
	if hashErr != nil {
		return model.JwtAuthRes{}, apperrors.ErrInvalidCredentials
		// w.WriteHeader(http.StatusBadRequest)
	}

	token, jwtErr := utils.JWTInit(userInfo.UserId, userInfo.Role)

	if jwtErr != nil {
		return model.JwtAuthRes{}, apperrors.ErrInternal
	}

	// fmt.Print(userInfo.UserId)

	jwtRes := model.JwtAuthRes{
		UserId:   userInfo.UserId,
		UserRole: userInfo.Role,
		JwtToken: token,
	}

	return jwtRes, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userId int, userUpdateInfo model.UserUpdate) (model.User, error) {
	if userId == 0 {
		return model.User{}, apperrors.ErrBadRequest
	}
	if userUpdateInfo.Email == nil && userUpdateInfo.Name == nil {
		return model.User{}, apperrors.ErrBadRequest
	}

	updatedUserRes, updateErr := s.repo.UpdateUser(ctx, userId, userUpdateInfo)

	if updateErr != nil {
		if errors.Is(updateErr, context.DeadlineExceeded) {
			return model.User{}, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(updateErr, context.Canceled) {
			return model.User{}, apperrors.ErrCanceled
		}

		return model.User{}, apperrors.ErrInternal
	}

	return updatedUserRes, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, apperrors.ErrBadRequest
	}

	deletedUserRes, deleteErr := s.repo.DeleteUser(ctx, userId)
	if deleteErr != nil {
		if errors.Is(deleteErr, sql.ErrNoRows) {
			return model.User{}, apperrors.ErrUserNotFound
		}
		return model.User{}, apperrors.ErrInternal
	}
	return deletedUserRes, nil
}

func (s *UserService) FetchUserById(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, apperrors.ErrBadRequest
	}

	userRes, userErr := s.repo.GetUserById(ctx, userId)
	if errors.Is(userErr, sql.ErrNoRows) {
		return model.User{}, apperrors.ErrUserNotFound
	}
	if userErr != nil {

		// fmt.Print(userErr)
		return model.User{}, apperrors.ErrInternal
	}

	return userRes, nil
}

func (s *UserService) FetchAllUsers(ctx context.Context) ([]model.User, error) {
	users, usersErr := s.repo.FetchAllUsers(ctx)
	if usersErr != nil {
		if errors.Is(usersErr, context.DeadlineExceeded) {
			return nil, apperrors.ErrDeadlineExceeded
		}
		if errors.Is(usersErr, context.Canceled) {
			return nil, apperrors.ErrCanceled
		}
		return nil, apperrors.ErrInternal
	}

	return users, nil
}
