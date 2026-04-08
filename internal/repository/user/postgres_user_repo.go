package user

import (
	"ai-task-processor/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, userDetail model.User) (int, error) {
	// println("in repo")
	// fmt.Printf("%+v", userDetail)
	createUserQuery := `INSERT INTO users(user_name, user_email, user_password) VALUES($1,$2,$3) RETURNING user_id;`

	row := r.db.QueryRowContext(ctx, createUserQuery, userDetail.Name, userDetail.Email, userDetail.Password)

	var newUserId int
	scanErr := row.Scan(&newUserId)
	if scanErr != nil {
		return 0, scanErr
	}
	return newUserId, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	// println(email)
	var user model.User
	fetchUserQuery := `SELECT user_id, user_name,user_email, user_password,user_role, created_at FROM users WHERE user_email = $1;`
	userInfo := r.db.QueryRowContext(ctx, fetchUserQuery, email)
	scanErr := userInfo.Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return model.User{}, sql.ErrNoRows
		}
		return model.User{}, scanErr
	}
	return user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userId int, userUpdateDetail model.UserUpdate) (model.User, error) {
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

	if len(args) <= 1 {
		return model.User{}, errors.New("No args passed in update")
	}

	updatedUser := r.db.QueryRowContext(ctx, updateQuery, args...)
	updateScanErr := updatedUser.Scan(&updatedUserRes.UserId, &updatedUserRes.Name, &updatedUserRes.Email, &updatedUserRes.CreatedAt)

	if updateScanErr != nil {
		return model.User{}, updateScanErr
	}

	return updatedUserRes, nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, userId int) (model.User, error) {
	var deletedUserRes model.User

	deleteQuery := `DELETE from users where user_id=$1 RETURNING user_id, user_name, user_email, created_at;`

	deletedUser := r.db.QueryRowContext(ctx, deleteQuery, userId)
	deleteScanErr := deletedUser.Scan(&deletedUserRes.UserId, &deletedUserRes.Name, &deletedUserRes.Email, &deletedUserRes.CreatedAt)
	if deleteScanErr != nil {
		if errors.Is(deleteScanErr, sql.ErrNoRows) {
			return model.User{}, sql.ErrNoRows
		}
		return model.User{}, deleteScanErr
	}

	return deletedUserRes, nil
}

func (r *PostgresUserRepository) GetUserById(ctx context.Context, userId int) (model.User, error) {
	var userRes model.User

	userFetchQuery := `SELECT user_id, user_name, user_email, user_role, created_at from users where user_id=$1;`
	// userFetchQuery := `SELECT pg_sleep(6)`

	userInfo := r.db.QueryRowContext(ctx, userFetchQuery, userId)
	userScanErr := userInfo.Scan(&userRes.UserId, &userRes.Name, &userRes.Email, &userRes.Role, &userRes.CreatedAt)

	if userScanErr != nil {
		if errors.Is(userScanErr, sql.ErrNoRows) {
			return model.User{}, sql.ErrNoRows
		}

		return model.User{}, userScanErr
	}

	return userRes, nil
}

func (r *PostgresUserRepository) FetchAllUsers(ctx context.Context) ([]model.User, error) {
	//initialise an array of type model.User
	var users []model.User

	//postgres query to fetch the whole data
	allUsersQuery := `SELECT user_id, user_name, user_email,user_role, created_at from users;`
	userRows, allUserErr := r.db.QueryContext(ctx, allUsersQuery)

	if allUserErr != nil {
		return nil, allUserErr
	}

	defer userRows.Close()
	//iterate over the rows and append each user one by one in users array
	for userRows.Next() {

		var user model.User
		scanErr := userRows.Scan(&user.UserId, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		users = append(users, user)
	}

	if rowsErr := userRows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	//return the users array

	return users, nil
}
