package repository

import (
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/model"
	"fmt"
)

// var db.Db *sql.DB = db.Db

func CreateUserRepo(userDetail model.User) (int, error) {
	// println("in repo")
	// fmt.Printf("%+v", userDetail)
	createUserQuery := `INSERT INTO users(user_name, user_email, user_password) VALUES($1,$2,$3) RETURNING user_id;`

	row := db.Db.QueryRow(createUserQuery, userDetail.Name, userDetail.Email, userDetail.Password)

	var newUserId int
	scanErr := row.Scan(&newUserId)
	if scanErr != nil {
		return 0, fmt.Errorf("internal db error: %w", scanErr)
	}
	return newUserId, nil
}

// get user info using email
func GetUserByEmail(email string) (model.User, error) {
	// println(email)
	var user model.User
	fetchUserQuery := `SELECT user_id, user_name,user_email, user_password, created_at FROM users WHERE user_email = $1;`
	userInfo := db.Db.QueryRow(fetchUserQuery, email)
	scanErr := userInfo.Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if scanErr != nil {
		return model.User{}, fmt.Errorf("fetch user failed:%w", scanErr)
	}
	return user, nil
}
