package handler

import (
	"ai-task-processor/internal/constants"
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

	utils.WriteJsonResponse(w, 201, true, "User Created ✅", userIdPayload)
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

	loginPayload, loginErr := service.LoginService(userLoginInfo)
	if loginErr != nil {
		utils.WriteJsonResponse(w, 401, false, loginErr.Error(), nil)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "Login Success", loginPayload)

}

// update user info
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, ok := r.Context().Value(constants.UserKey).(int)
	fmt.Print(userId)
	if !ok {
		utils.WriteJsonResponse(w, 401, false, "Unauthorized1", nil)
		return
	}
	var updateUserInfo model.UserUpdate

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&updateUserInfo)
	if decodeErr != nil {
		utils.WriteJsonResponse(w, 500, false, "Update decode failed", nil)
		return
	}

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

	userId := r.Context().Value(constants.UserKey).(int)

	deletedUserRes, deleteErr := service.DeleteUserService(userId)
	if deleteErr != nil {
		utils.WriteJsonResponse(w, 500, false, deleteErr.Error(), nil)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "User Deleted", deletedUserRes)

}
