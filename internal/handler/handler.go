package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"ai-task-processor/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
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
