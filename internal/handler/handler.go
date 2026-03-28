package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"encoding/json"
	"fmt"
	"log"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service.CreateUserService(newUser)
	w.WriteHeader(http.StatusCreated)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userLoginInfo model.UserLogin
	decoder := json.NewDecoder(r.Body)
	loginErr := decoder.Decode(&userLoginInfo)
	if loginErr != nil {
		log.Printf("Error while logging in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	service.LoginService(userLoginInfo)

}
