package dbHelper

import (
	"Todo-Server/database"
	"Todo-Server/models"
)

func GetUsersCount() (int, error) {
	var count int
	query := `SELECT COUNT(*)
        FROM users
        WHERE archived_at IS NULL;`

	err := database.Todo.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FetchAllUsers(limit, offset int) ([]models.User, error) {
	var users []models.User

	SQL := `SELECT id, name, email, created_at
        FROM users
        WHERE archived_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2;`

	err := database.Todo.Select(&users, SQL, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func FetchUsersCount() (int, error) {
	var count int

	SQL := `SELECT COUNT(*)
		FROM users
		WHERE archived_at IS NULL;`

	err := database.Todo.Get(&count, SQL)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FetchAllTodos(limit, offset int) ([]models.Todo, error) {
	var todos []models.Todo

	SQL := `SELECT id, user_id, name , is_completed, created_at
		FROM todos
		WHERE archived_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;`

	err := database.Todo.Select(&todos, SQL, limit, offset)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func FetchTodosCount() (int, error) {
	var count int

	SQL := `SELECT COUNT(*)
		FROM todos
		WHERE archived_at IS NULL;`

	err := database.Todo.Get(&count, SQL)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func SuspendUser(userID string) error {
	SQL := `UPDATE users
	        SET suspended_at = NOW()
	        WHERE id = $1 AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, userID)
	return err
}

func IsUserSuspended(userID string) (bool, error) {
	SQL := `SELECT suspended_at IS NOT NULL
	        FROM users
	        WHERE id = $1`

	var suspended bool
	err := database.Todo.Get(&suspended, SQL, userID)
	return suspended, err
}
