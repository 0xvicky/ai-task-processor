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
	if errors.Is(fetchErr, sql.ErrNoRows) {
		return model.JwtAuthRes{}, fmt.Errorf("invalid credentials")
	}
	if fetchErr != nil {
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

func MeService(userId int) (model.User, error) {
	if userId == 0 {
		return model.User{}, fmt.Errorf("Invalid User")
	}

	userRes, userErr := repository.MeRepo(userId)
	if errors.Is(userErr, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("User not found")
	}
	if userErr != nil {
		fmt.Print(userErr)
		return model.User{}, fmt.Errorf("Invalid User")
	}

	return userRes, nil

}

func FetchAllUsersService() ([]model.User, error) {
	users, usersErr := repository.FetchAllUsersRepo()
	if usersErr != nil {
		return nil, fmt.Errorf("fetch all users failed")
	}

	return users, nil
}
