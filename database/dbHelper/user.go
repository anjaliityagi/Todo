package dbHelper

import (
	"Todo-Server/database"
	"Todo-Server/models"

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

	//var todo models.Todos
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
