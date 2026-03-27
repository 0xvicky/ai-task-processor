package repository

import (
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/model"
	"fmt"
	"log"
)

// var db.Db *sql.DB = db.Db

func CreateUserRepo(userDetail model.User) {
	println("in repo")
	fmt.Printf("%+v", userDetail)
	//create table if not exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS users(
	user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_name TEXT NOT NULL,
	user_email TEXT NOT NULL UNIQUE,
	user_password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, createTableErr := db.Db.Exec(createTableQuery)
	if createTableErr != nil {
		log.Fatalf("Error occured while creating user:%v", createTableErr)
	}

	fmt.Println("User table will be there !!!")

	createUserQuery := `INSERT INTO users(user_name, user_email, user_password) VALUES($1,$2,$3);`

	res, createUserErr := db.Db.Exec(createUserQuery, userDetail.Name, userDetail.Email, userDetail.Password)

	if createUserErr != nil {
		log.Fatalf("Error occured while creating user:%v", createUserErr)
	}
	fmt.Printf("%v", res)

}
