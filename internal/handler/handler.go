package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"ai-task-processor/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// default route for "/"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ai Task Processor running on port 6969")
}

func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Health is OK")
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newUser model.User
	// println(r.Body)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newUser)
	if err != nil {
		// Handle errors (e.g., malformed JSON, wrong type)
		utils.WriteJsonResponse(w, 400, false, err.Error(), nil)
		return
	}

	userIdPayload, userCreateErr := service.CreateUserService(newUser)
	if userCreateErr != nil {
		utils.WriteJsonResponse(w, 500, false, userCreateErr.Error(), nil)
		return
	}
	// w.WriteHeader(http.StatusCreated)
	payload := struct {
		UserId int `json:"user_id"`
	}{
		UserId: userIdPayload,
	}
	utils.WriteJsonResponse(w, 201, true, "User Created ✅", payload)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userLoginInfo model.UserLogin
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&userLoginInfo)
	if decodeErr != nil {
		utils.WriteJsonResponse(w, 400, false, decodeErr.Error(), nil)
		return
	}

	loginToken, loginErr := service.LoginService(userLoginInfo)
	if loginErr != nil {
		utils.WriteJsonResponse(w, 401, false, loginErr.Error(), nil)
		return
	}

	loginPayload := struct {
		Token string `json:"token"`
	}{
		Token: loginToken,
	}

	utils.WriteJsonResponse(w, 200, true, "Login Success", loginPayload)

}

// update user info
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var updateUserInfo model.UserUpdate
	var userId int

	var req struct {
		UserId int `json:"userid"`
		model.UserUpdate
	}

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&req)
	if decodeErr != nil {
		utils.WriteJsonResponse(w, 500, false, "Update decode failed", nil)
		return
	}
	userId = req.UserId
	updateUserInfo = req.UserUpdate

	userUpdateRes, updateErr := service.UpdateService(updateUserInfo, userId)

	if updateErr != nil {
		utils.WriteJsonResponse(w, 500, false, updateErr.Error(), nil)
		return
	}

	updatePayload := struct {
		UserId    int       `json:"user_id"`
		Name      string    `json:"user_name"`
		Email     string    `json:"user_email"`
		CreatedAt time.Time `json:"created_at"`
	}{
		UserId:    userUpdateRes.UserId,
		Name:      userUpdateRes.Name,
		Email:     userUpdateRes.Email,
		CreatedAt: userUpdateRes.CreatedAt,
	}

	utils.WriteJsonResponse(w, 200, true, "Update Success", updatePayload)

}

// delete user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// var userId int `json:"userid"`
	var req struct {
		UserId int `json:"userid"`
	}
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&req)
	if decodeErr != nil {
		utils.WriteJsonResponse(w, 400, false, decodeErr.Error(), nil)
	}

	deletedUserRes, deleteErr := service.DeleteUserService(req.UserId)
	if deleteErr != nil {
		utils.WriteJsonResponse(w, 500, false, deleteErr.Error(), nil)
	}

	utils.WriteJsonResponse(w, 200, true, "User Deleted", deletedUserRes)

}
