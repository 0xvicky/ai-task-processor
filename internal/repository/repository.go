package repository

import (
	"ai-task-processor/internal/db"
	"ai-task-processor/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// var db.Db *sql.DB = db.Db

func CreateUserRepo(userDetail model.User) (int, error) {
	// println("in repo")
	// fmt.Printf("%+v", userDetail)
	createUserQuery := `INSERT INTO users(user_name, user_email, user_password, user_role) VALUES($1,$2,$3, $4) RETURNING user_id;`

	row := db.Db.QueryRow(createUserQuery, userDetail.Name, userDetail.Email, userDetail.Password, userDetail.Role)

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
	fetchUserQuery := `SELECT user_id, user_name,user_email, user_password,user_role, created_at FROM users WHERE user_email = $1;`
	userInfo := db.Db.QueryRow(fetchUserQuery, email)
	scanErr := userInfo.Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if scanErr != nil {
		if errors.Is(scanErr, context.DeadlineExceeded) {
			return model.User{}, context.DeadlineExceeded
		}
		if errors.Is(scanErr, context.Canceled) || strings.Contains(scanErr.Error(), "canceling statement") {
			return model.User{}, context.Canceled
		}
		return model.User{}, fmt.Errorf("fetch user failed:%w", scanErr)
	}
	return user, nil
}

func UpdateUserRepo(ctx context.Context, userId int, userUpdateDetail model.UserUpdate) (model.User, error) {
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
		if errors.Is(updateScanErr, context.DeadlineExceeded) {
			return model.User{}, context.DeadlineExceeded
		}
		if errors.Is(updateScanErr, context.Canceled) || strings.Contains(updateScanErr.Error(), "canceling statements") {
			return model.User{}, context.Canceled
		}
		return model.User{}, fmt.Errorf("Update Scanner failed ! %w", updateScanErr)
	}

	return updatedUserRes, nil
}

func DeleteUserRepo(ctx context.Context, userId int) (model.User, error) {
	var deletedUserRes model.User

	deleteQuery := `DELETE from users where user_id=$1 RETURNING user_id, user_name, user_email, created_at;`

	deletedUser := db.Db.QueryRowContext(ctx, deleteQuery, userId)
	deleteScanErr := deletedUser.Scan(&deletedUserRes.UserId, &deletedUserRes.Name, &deletedUserRes.Email, &deletedUserRes.CreatedAt)
	if deleteScanErr != nil {
		if errors.Is(deleteScanErr, context.DeadlineExceeded) {
			return model.User{}, context.DeadlineExceeded
		}
		if errors.Is(deleteScanErr, context.Canceled) || strings.Contains(deleteScanErr.Error(), "canceling statement") {
			return model.User{}, context.Canceled
		}
		return model.User{}, deleteScanErr
	}

	return deletedUserRes, nil
}

func MeRepo(ctx context.Context, userId int) (model.User, error) {
	var userRes model.User

	userFetchQuery := `SELECT user_id, user_name, user_email, created_atfrom users where user_id=$1;`
	// userFetchQuery := `SELECT pg_sleep(6)`

	userInfo := db.Db.QueryRowContext(ctx, userFetchQuery, userId)
	userScanErr := userInfo.Scan(&userRes.UserId, &userRes.Name, &userRes.Email, &userRes.CreatedAt)

	if errors.Is(userScanErr, context.DeadlineExceeded) {
		return model.User{}, context.DeadlineExceeded
	}
	if errors.Is(userScanErr, context.Canceled) || strings.Contains(userScanErr.Error(), "canceling statement") {
		return model.User{}, context.Canceled
	}

	if errors.Is(userScanErr, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("User not found: %w", sql.ErrNoRows)
	}

	if userScanErr != nil {
		return model.User{}, fmt.Errorf("User Fetch scanner failed %w", userScanErr)
	}

	return userRes, nil
}

func FetchAllUsersRepo(ctx context.Context) ([]model.User, error) {

	//initialise an array of type model.User
	var users []model.User

	//postgres query to fetch the whole data
	allUsersQuery := `SELECT user_id, user_name, user_email, created_at from users;`
	userRows, allUserErr := db.Db.QueryContext(ctx, allUsersQuery)

	if allUserErr != nil {
		if errors.Is(allUserErr, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		if errors.Is(allUserErr, context.Canceled) || strings.Contains(allUserErr.Error(), "canceling statement") {
			return nil, context.Canceled
		}
		return nil, fmt.Errorf("Fetch all users error")
	}

	defer userRows.Close()
	//iterate over the rows and append each user one by one in users array
	for userRows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		var user model.User
		scanErr := userRows.Scan(&user.UserId, &user.Name, &user.Email, &user.CreatedAt)
		if scanErr != nil {
			return nil, fmt.Errorf("all users scan failed %w", scanErr)
		}
		users = append(users, user)
	}

	if rowsErr := userRows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("rows iteration error %w", rowsErr)
	}
	//return the users array
	return users, nil
}
