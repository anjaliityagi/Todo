package dbHelper

import (
	"Todo-Server/database"
	"Todo-Server/models"
	//"Todo-Server/utils"
	"database/sql"
	"errors"
	"time"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT count(*) > 0
			FROM users
			WHERE email = TRIM(LOWER($1))
			  AND archived_at IS NULL;`

	var exists bool
	err := database.Todo.Get(&exists, SQL, email)
	return exists, err
}

func CreateUser(name, email, password string) (string, error) {
	SQL := `INSERT INTO users(name, email, password)
			VALUES ($1, TRIM(LOWER($2)), $3)
			RETURNING id;`

	var userID string
	err := database.Todo.Get(&userID, SQL, name, email, password)
	return userID, err
}

func CreateTodo(userID, name, description string, expiringAt time.Time) (models.Todos, error) {
	SQL := `INSERT INTO todos (user_id,name,description,expiring_at) 
			VALUES ($1,$2,$3,$4)
			RETURNING id,user_id,name,description,complete,expiring_at,created_at;`

	var todo models.Todos
	err := database.Todo.Get(&todo, SQL, userID, name, description, expiringAt)
	return todo, err
}

func GetTodos(userID, search, date, complete string) ([]models.Todos, error) {
	SQL := `
		SELECT id,user_id,name,description,complete,expiring_at,created_at
		FROM todos
		WHERE user_id = $1
		AND archived_at IS NULL
		AND ($2 = '' OR complete = $2::boolean)
		AND ($3 = '' OR expiring_at <= $3::timestamptz)
		AND ($4 = '' OR name ILIKE '%' || $4 || '%')
		ORDER BY expiring_at;
	`

	var todos []models.Todos
	err := database.Todo.Select(&todos, SQL, userID, complete, date, search)
	return todos, err
}

func GetTodoByID(todoID, userID string) (*models.Todos, error) {
	SQL := `SELECT id,user_id,name,description,complete,expiring_at,created_at
			FROM todos
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var todo models.Todos

	err := database.Todo.Get(&todo, SQL, todoID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	return &todo, nil
}

func DeleteTodoById(userID, todoID string) error {
	SQL := `UPDATE todos
			SET archived_at = NOW()
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	result, err := database.Todo.Exec(SQL, todoID, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateTodoById(
	name string,
	description string,
	complete bool,
	expiringAt time.Time,
	todoID string,
	userID string,
) error {

	SQL := `UPDATE todos 
			SET name=$1, description=$2, complete=$3, expiring_at=$4
			WHERE id=$5 AND user_id=$6 AND archived_at IS NULL`

	result, err := database.Todo.Exec(
		SQL,
		name,
		description,
		complete,
		expiringAt,
		todoID,
		userID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

//func CreateUserSession(userID string) (string, error) {
//	SQL := `INSERT INTO user_session(user_id)
//			VALUES ($1) RETURNING id;`
//	var sessionID string
//	err := database.Todo.Get(&sessionID, SQL, userID)
//	if err != nil {
//		return "", err
//	}
//	return sessionID, nil
//}
//func GetUserByEmail(email, password string) (string, error) {
//	SQL := `
//		SELECT id, password
//		FROM users
//		WHERE email = $1 AND archived_at IS NULL;
//	`
//
//	var user models.UserAuth
//	err := database.Todo.Get(&user, SQL, email)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return "", errors.New("no user exist")
//		}
//		return "", err
//	}
//	if err := utils.CheckPassword(user.Password, password); err != nil {
//		return "", errors.New("invalid credentials")
//	}
//	return user.ID, nil
//}
