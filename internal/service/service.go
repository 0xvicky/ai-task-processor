package service

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(newUser model.User) (model.JwtAuthRes, error) {
	_, fetchErr := repository.GetUserByEmail(newUser.Email)

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
	userId, creationErr := repository.CreateUserRepo(newUser)
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

func LoginService(userLoginInfo model.UserLogin) (model.JwtAuthRes, error) {
	//check if user exist or not, if exist fetch the user details
	userInfo, fetchErr := repository.GetUserByEmail(userLoginInfo.Email)
	if fetchErr != nil {

		if errors.Is(fetchErr, sql.ErrNoRows) {
			return model.JwtAuthRes{}, sql.ErrNoRows
		}
		if errors.Is(fetchErr, context.DeadlineExceeded) {
			return model.JwtAuthRes{}, context.DeadlineExceeded
		}
		if errors.Is(fetchErr, context.Canceled) || strings.Contains(fetchErr.Error(), "canceling statement") {
			return model.JwtAuthRes{}, context.Canceled
		}
		return model.JwtAuthRes{}, fmt.Errorf("fetch user failed:%w", fetchErr)
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

func UpdateService(ctx context.Context, userId int, userUpdateInfo model.UserUpdate) (model.User, error) {

	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User Id")
	}
	if userUpdateInfo.Email == nil && userUpdateInfo.Name == nil {
		return model.User{}, fmt.Errorf("Both fields are empty !")
	}

	updatedUserRes, updateErr := repository.UpdateUserRepo(ctx, userId, userUpdateInfo)

	if errors.Is(updateErr, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("No user found:%w", updateErr)
	}
	if updateErr != nil {
		return model.User{}, fmt.Errorf("Update user failed:%w", updateErr)
	}

	return updatedUserRes, nil
}

func DeleteUserService(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User Id")
	}

	deletedUserRes, deleteErr := repository.DeleteUserRepo(ctx, userId)
	if errors.Is(deleteErr, sql.ErrNoRows) {
		return model.User{UserId: userId}, nil
	}
	if deleteErr != nil {
		return model.User{}, deleteErr
	}
	return deletedUserRes, nil

}

func MeService(ctx context.Context, userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User")
	}

	userRes, userErr := repository.MeRepo(ctx, userId)

	if errors.Is(userErr, context.DeadlineExceeded) {
		return model.User{}, userErr
	}
	if errors.Is(userErr, sql.ErrNoRows) {
		return model.User{}, sql.ErrNoRows
	}
	if userErr != nil {
		// fmt.Print(userErr)
		return model.User{}, userErr
	}

	return userRes, nil

}

func FetchAllUsersService(ctx context.Context) ([]model.User, error) {
	users, usersErr := repository.FetchAllUsersRepo(ctx)
	if usersErr != nil {
		return nil, usersErr
	}

	return users, nil
}
