package repository

import (
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/model"
	"fmt"
	"strings"
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

func UpdateUserRepo(userUpdateDetail model.UserUpdate, userId int) (model.User, error) {
	var updatedUserRes model.User
	updateQuery := `UPDATE users SET `
	args := []any{}
	i := 1

	if userUpdateDetail.Email != nil {
		updateQuery += fmt.Sprintf(`user_email = $%d, `, i)
		args = append(args, *userUpdateDetail.Email)
		i++
	}
	if userUpdateDetail.Name != nil {
		updateQuery += fmt.Sprintf(`user_name=$%d, `, i)
		args = append(args, *userUpdateDetail.Name)
		i++
	}
	updateQuery = strings.TrimSuffix(updateQuery, ", ")
	updateQuery += fmt.Sprintf(` WHERE user_id=$%d RETURNING user_id, user_name, user_email, created_at`, i)
	args = append(args, userId)

	if len(args) == 0 {
		return model.User{}, fmt.Errorf("No args passed in update")
	}
	updatedUser := db.Db.QueryRow(updateQuery, args...)
	updateScanErr := updatedUser.Scan(&updatedUserRes.UserId, &updatedUserRes.Name, &updatedUserRes.Email, &updatedUserRes.CreatedAt)

	if updateScanErr != nil {
		return model.User{}, fmt.Errorf("Update Scanner failed ! %w", updateScanErr)
	}

	return updatedUserRes, nil
}

func DeleteUserRepo(userId int) (model.User, error) {
	var deletedUserRes model.User

	deleteQuery := `DELETE from users where user_id=$1 RETURNING user_id, user_name, user_email, created_at;`

	deletedUser := db.Db.QueryRow(deleteQuery, userId)
	deleteScanErr := deletedUser.Scan(&deletedUserRes.UserId, &deletedUserRes.Name, &deletedUserRes.Email, &deletedUserRes.CreatedAt)
	if deleteScanErr != nil {
		return model.User{}, fmt.Errorf("Delete Scanner Failed:%w", deleteScanErr)
	}

	return deletedUserRes, nil
}

func MeRepo(userId int) (model.User, error) {
	var userRes model.User

	userFetchQuery := `SELECT user_id, user_name, user_email, created_at from users where user_id=$1;`

	userInfo := db.Db.QueryRow(userFetchQuery, userId)
	userScanErr := userInfo.Scan(&userRes.UserId, &userRes.Name, &userRes.Email, &userRes.CreatedAt)

	if userScanErr != nil {
		return model.User{}, fmt.Errorf("User Fetch scanner failed%w", userScanErr)
	}

	return userRes, nil

}
