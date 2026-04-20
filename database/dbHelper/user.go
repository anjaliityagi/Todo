package dbHelper

import (
	"Todo-Server/database"
	"Todo-Server/models"
	"fmt"
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

func CreateUser(name, email, password string) error {
	SQL := `INSERT INTO users(name, email, password)
			VALUES ($1, TRIM(LOWER($2)), $3)`

	_, err := database.Todo.Exec(SQL, name, email, password)
	return err
}

func CreateTodo(userID, name, description string, expiringAt time.Time) error {
	SQL := `INSERT INTO todos (user_id,name,description,expiring_at)
			VALUES ($1,$2,$3,$4)`

	_, err := database.Todo.Exec(SQL, userID, name, description, expiringAt)
	return err
}

func CreateUserSession(userID string) (string, error) {
	SQL := `INSERT INTO user_session(user_id)
			VALUES ($1) RETURNING id;`
	var sessionID string
	err := database.Todo.Get(&sessionID, SQL, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
func GetUserByEmail(email string) (*models.UserAuth, error) {
	SQL := `
		SELECT id, password
		FROM users
		WHERE email = lower(trim($1)) AND archived_at IS NULL;
	`

	var user models.UserAuth
	err := database.Todo.Get(&user, SQL, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteSessionByToken(token string) error {
	SQL := `UPDATE user_session
			SET archived_at = NOW()
			WHERE id = $1 AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, token)
	if err != nil {
		return err
	}

	return nil
}

func GetUserIDBySession(token string) (string, error) {
	SQL := `SELECT user_id FROM user_session WHERE id = $1 and archived_at IS NULL `

	var userId string

	fmt.Println(userId)

	err := database.Todo.Get(&userId, SQL, token)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func UpdateTodo(id, userID, name string, description string, expiringAt time.Time, complete bool) error {

	SQL := `UPDATE todos
SET name = $3,
    description = $4,
    expiring_at = $5,
    complete=$6
WHERE id = $1
  AND user_id = $2
  AND archived_at IS NULL; `

	_, err := database.Todo.Exec(SQL, id, userID, name, description, expiringAt, complete)
	return err

}

func DeleteTodo(id, userId string) error {
	SQL := `UPDATE todos
SET archived_at = NOW()
where id = $1
  AND user_id = $2`

	_, err := database.Todo.Exec(SQL, id, userId)
	if err != nil {
		return err
	}

	return nil
}

func GetTodoById(todoId string, userId string) (models.Todo, error) {
	var todo models.Todo
	SQL := `
SELECT name, description, complete, expiring_at, created_at
FROM todos
WHERE id = $1
  AND user_id = $2 `
	err := database.Todo.Get(&todo, SQL, todoId, userId)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}
func GetTodos(userID, search, date, status string) ([]models.Todo, error) {

	SQL := `SELECT id,
	       user_id,
	       name,
	       description,
	       complete,
	       expiring_at,
	       created_at
	FROM todos
	WHERE user_id = $1
	  AND archived_at IS NULL`

	args := []interface{}{userID}
	i := 2
	if status != "" {
		if status == "completed" {
			SQL += fmt.Sprintf(" AND complete = $%d", i)
			args = append(args, true)
			i++

		} else if status == "pending" {
			SQL += " AND complete = false AND expiring_at >= NOW()"

		} else if status == "expired" {
			SQL += " AND complete = false AND expiring_at < NOW()"

		} else {
			return nil, fmt.Errorf("invalid status")
		}
	}
	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format (use YYYY-MM-DD)")
		}
		SQL += fmt.Sprintf(" AND expiring_at <= $%d", i)
		args = append(args, t)
		i++
	}

	if search != "" {
		SQL += fmt.Sprintf(" AND name ILIKE $%d", i)
		args = append(args, "%"+search+"%")
		i++
	}

	SQL += " ORDER BY expiring_at"

	var todos []models.Todo

	err := database.Todo.Select(&todos, SQL, args...)
	if err != nil {
		return nil, err
	}
	return todos, nil
}
