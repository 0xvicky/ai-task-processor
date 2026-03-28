package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
	userDet := repository.GetUserByEmail(newUser.Email)

	if newUser.Email == userDet.Email {
		log.Printf("User already exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//get the password out of body
	plainPass := newUser.Password
	//hash the password
	hashPass, hashErr := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)

	if hashErr != nil {
		log.Fatal("Fault during hashing pass:", hashErr)
	}
	//save the string(hash) back to struct
	newUser.Password = string(hashPass)
	//passing the user details to the db-repository
	repository.CreateUserRepo(newUser)
	w.WriteHeader(http.StatusCreated)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userLoginInfo := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	loginErr := decoder.Decode(&userLoginInfo)
	if loginErr != nil {
		log.Printf("Error while logging in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check if user exist or not, if exist fetch the user details
	userInfo := repository.GetUserByEmail(userLoginInfo.Email)
	fmt.Printf("%+v", userInfo)
	if userInfo.Email == "" {
		log.Printf("User Doesn't exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if exist, then hash the password and compare with stored hash pass

	hashErr := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(userLoginInfo.Password))
	if hashErr != nil {
		log.Printf("Wrong password while loggin")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
