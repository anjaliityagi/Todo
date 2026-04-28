package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CreateTodo struct {
	Name        string    `json:"name" binding:"required,max=30"`
	Description string    `json:"description" binding:"required,max=200"`
	ExpiringAt  time.Time `json:"expiringAt" binding:"required"`
}

type RegisterUser struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}
type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type UserAuth struct {
	ID       string `db:"id"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type UpdateTodo struct {
	Name        *string    `json:"name" binding:"omitempty,max=30"`
	Description *string    `json:"description" binding:"omitempty,max=200"`
	Complete    *bool      `json:"isCompleted" `
	ExpiringAt  *time.Time `json:"expiringAt"`
}

type Todo struct {
	TodoID      string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"userId"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsCompleted bool      `db:"is_completed" json:"isCompleted"`
	ExpiringAt  time.Time `db:"expiring_at" json:"expiringAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type User struct {
	ID        string `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Email     string `db:"email" json:"email"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}
