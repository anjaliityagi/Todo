package dbHelper

import (
	"Todo-Server/database"
	"Todo-Server/models"
	"fmt"
	"strings"
	"time"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT count(*) > 0
			FROM users
			WHERE email = TRIM(LOWER($1))
			  AND archived_at IS NULL`

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

func CreateTodo(userID, name, description string, expiringAt time.Time) (string, error) {
	var todoID string
	SQL := `INSERT INTO todos (user_id,name,description,expiring_at)
			VALUES ($1,$2,$3,$4)
			RETURNING id`

	err := database.Todo.Get(&todoID, SQL, userID, name, description, expiringAt)

	if err != nil {
		return "", err
	}
	return todoID, nil
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
	SQL := `SELECT u.id, u.password, COALESCE(ur.role, 'user') AS role
FROM users u
LEFT JOIN user_roles ur 
  ON u.id = ur.user_id AND ur.archived_at IS NULL
WHERE u.email = lower(trim($1)) 
  AND u.archived_at IS NULL
  AND u.suspended_at IS NULL;`

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

//func GetUserIDBySession(token string) (string, error) {
//	SQL := `SELECT user_id
//FROM user_session
//WHERE id = $1
//and archived_at IS NULL `
//
//	var userId string
//
//	fmt.Println(userId)
//
//	err := database.Todo.Get(&userId, SQL, token)
//	if err != nil {
//		return "", err
//	}
//	return userId, nil
//}

func DeleteTodo(id, userId string) error {
	SQL := `UPDATE todos
SET archived_at = NOW()
where id = $1
  AND user_id = $2
  AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, id, userId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTodoByID(id, userID string, updateTodo models.UpdateTodo) error {

	SQL := "UPDATE todos SET "
	args := []interface{}{}

	if updateTodo.Name != nil {
		args = append(args, *updateTodo.Name)

		SQL += fmt.Sprintf("name = $%d, ", len(args))
	}
	if updateTodo.Description != nil {
		args = append(args, *updateTodo.Description)
		SQL += fmt.Sprintf("description = $%d, ", len(args))
	}
	if updateTodo.ExpiringAt != nil {
		args = append(args, *updateTodo.ExpiringAt)
		SQL += fmt.Sprintf("expiring_at = $%d, ", len(args))
	}
	if updateTodo.Complete != nil {
		args = append(args, *updateTodo.Complete)
		SQL += fmt.Sprintf("is_completed = $%d, ", len(args))
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	SQL = strings.TrimSuffix(SQL, ", ")
	IdIndex := len(args) + 1
	UserIdIndex := len(args) + 2
	SQL += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d AND archived_at IS NULL", IdIndex, UserIdIndex)
	args = append(args, id, userID)

	_, err := database.Todo.Exec(SQL, args...)
	return err
}

func FetchTodoById(todoId string, userId string) (models.Todo, error) {
	var todo models.Todo
	SQL := `SELECT id,user_id,name, description, is_completed, expiring_at, created_at
FROM todos
WHERE id = $1
  AND user_id = $2`

	err := database.Todo.Get(&todo, SQL, todoId, userId)

	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func FetchTodos(userID, search, date, status string, limit int, page int, offset int) ([]models.Todo, error) {

	SQL := `SELECT id,
	       user_id,
	       name,
	       description,
	       is_completed,
	       expiring_at,
	       created_at
	FROM todos
	WHERE user_id = $1
	  AND archived_at IS NULL
	  AND (
		$2 = '' 
		OR ($2 = 'completed' AND is_completed = true)
		OR ($2 = 'pending' AND is_completed = false AND expiring_at >= NOW())
		OR ($2 = 'expired' AND is_completed = false AND expiring_at < NOW())
	  )
	  AND (
		$3 = '' OR expiring_at <= $3::TIMESTAMPTZ
	  )
	  AND (
		$4 = '' OR name ILIKE '%' || $4 || '%'
	  )
	ORDER BY expiring_at
	LIMIT $5 OFFSET $6
	`

	todos := make([]models.Todo, 0)

	err := database.Todo.Select(&todos, SQL, userID, status, date, search, limit, offset)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func FetchTotalTodoCount(userId string) (int, error) {
	var total int
	SQL := `select Count(*) from todos where user_id=$1 and archived_at is null `
	err := database.Todo.Get(&total, SQL, userId)
	if err != nil {
		return 0, err
	}
	return total, err
}

func IsSessionActive(sessionID string) (bool, error) {
	SQL := `SELECT count(*) > 0
	        FROM user_session
	        WHERE id = $1
	          AND archived_at IS NULL`

	var exists bool
	err := database.Todo.Get(&exists, SQL, sessionID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
