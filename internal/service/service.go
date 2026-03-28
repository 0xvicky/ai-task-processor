package service

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
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

	return "DUMMYTOKEN", nil

}
