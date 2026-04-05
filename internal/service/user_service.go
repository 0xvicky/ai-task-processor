package service

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return model.JwtAuthRes{}, fmt.Errorf("user already exist")
	}
	if !errors.Is(fetchErr, sql.ErrNoRows) {
		return model.JwtAuthRes{}, fmt.Errorf("internal db error")
	}

	//get the password out of body
	plainPass := newUser.Password
	//hash the password
	hashPass, hashErr := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)

	if hashErr != nil {
		return model.JwtAuthRes{}, fmt.Errorf("password hashing failed")
	}
	//save the string(hash) back to struct
	newUser.Password = string(hashPass)
	//passing the user details to the db-repository
	userId, creationErr := s.repo.CreateUser(ctx, newUser)
	if creationErr != nil {
		return model.JwtAuthRes{}, fmt.Errorf("internal db error")
	}

	token, jwtErr := utils.JWTInit(userId, newUser.Role)

	if jwtErr != nil {
		return model.JwtAuthRes{}, fmt.Errorf("Error while generating jwt:%w", jwtErr)
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

		return model.JwtAuthRes{}, fmt.Errorf("invalid credentials:%w", fetchErr)
	}
	// //if exist, then hash the password and compare with stored hash pass

	hashErr := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(userLoginInfo.Password))
	if hashErr != nil {
		return model.JwtAuthRes{}, fmt.Errorf("invalid credentails")
		// w.WriteHeader(http.StatusBadRequest)
	}

	token, jwtErr := utils.JWTInit(userInfo.UserId, userInfo.Role)

	if jwtErr != nil {
		return model.JwtAuthRes{}, fmt.Errorf("Error while generating jwt:%w", jwtErr)
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
		return model.User{}, fmt.Errorf("invalid User Id")
	}
	if userUpdateInfo.Email == nil && userUpdateInfo.Name == nil {
		return model.User{}, fmt.Errorf("Both fields are empty !")
	}

	updatedUserRes, updateErr := s.repo.UpdateUser(ctx, userId, userUpdateInfo)

	if updateErr != nil {

		return model.User{}, fmt.Errorf("Update user failed:%w", updateErr)
	}

	return updatedUserRes, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User Id")
	}

	deletedUserRes, deleteErr := s.repo.DeleteUser(ctx, userId)
	if deleteErr != nil {

		return model.User{}, deleteErr
	}
	return deletedUserRes, nil
}

func (s *UserService) FetchUserById(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User")
	}

	userRes, userErr := s.repo.GetUserById(ctx, userId)

	if userErr != nil {

		// fmt.Print(userErr)
		return model.User{}, userErr
	}

	return userRes, nil
}

func (s *UserService) FetchAllUsers(ctx context.Context) ([]model.User, error) {
	users, usersErr := s.repo.FetchAllUsers(ctx)
	if usersErr != nil {

		return nil, usersErr
	}

	return users, nil
}
