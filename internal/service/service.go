package service

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
	"ai-task-processor/internal/utils"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(newUser model.User) (int, error) {
	_, fetchErr := repository.GetUserByEmail(newUser.Email)

	if fetchErr == nil {
		return 0, fmt.Errorf("user already exist")
	}
	if !errors.Is(fetchErr, sql.ErrNoRows) {
		return 0, fmt.Errorf("internal db error")
	}

	//get the password out of body
	plainPass := newUser.Password
	//hash the password
	hashPass, hashErr := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)

	if hashErr != nil {
		return 0, fmt.Errorf("password hashing failed")
	}
	//save the string(hash) back to struct
	newUser.Password = string(hashPass)
	//passing the user details to the db-repository
	userId, creationErr := repository.CreateUserRepo(newUser)
	if creationErr != nil {
		return 0, fmt.Errorf("internal db error")
	}

	return userId, nil
}

func LoginService(userLoginInfo model.UserLogin) (string, error) {
	//check if user exist or not, if exist fetch the user details
	userInfo, fetchErr := repository.GetUserByEmail(userLoginInfo.Email)
	if errors.Is(fetchErr, sql.ErrNoRows) {
		return "", fmt.Errorf("invalid credentials")
	}
	if fetchErr != nil {
		return "", fmt.Errorf("fetch user failed:%w", fetchErr)
	}
	// //if exist, then hash the password and compare with stored hash pass

	hashErr := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(userLoginInfo.Password))
	if hashErr != nil {
		return "", fmt.Errorf("invalid credentails")
		// w.WriteHeader(http.StatusBadRequest)
	}

	//generate jwt token using userId and email
	jwtInfo := model.JWTModel{
		UserId: &userInfo.UserId,
		Email:  &userInfo.Email,
	}

	token, jwtErr := utils.JWTInit(jwtInfo)

	if jwtErr != nil {
		return "", fmt.Errorf("Error while generating jwt:%w", jwtErr)
	}

	return token, nil

}

func UpdateService(userUpdateInfo model.UserUpdate, userId int) (model.User, error) {

	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User Id")
	}
	if userUpdateInfo.Email == nil && userUpdateInfo.Name == nil {
		return model.User{}, fmt.Errorf("Both fields are empty !")
	}

	updatedUserRes, updateErr := repository.UpdateUserRepo(userUpdateInfo, userId)

	if errors.Is(updateErr, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("No user found:%w", updateErr)
	}
	if updateErr != nil {
		return model.User{}, fmt.Errorf("Update user failed:%w", updateErr)
	}

	return updatedUserRes, nil
}

func DeleteUserService(userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User Id")
	}

	deletedUserRes, deleteErr := repository.DeleteUserRepo(userId)
	if errors.Is(deleteErr, sql.ErrNoRows) {
		return model.User{UserId: userId}, nil
	}
	if deleteErr != nil {
		return model.User{}, fmt.Errorf("Delete User Failed:%w", deleteErr)
	}
	return deletedUserRes, nil

}
