package models

import "time"

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
}

type UpdateTodo struct {
	Name        *string    `json:"name" binding:"omitempty,max=30"`
	Description *string    `json:"description" binding:"omitempty,max=200"`
	Complete    *bool      `json:"complete" `
	ExpiringAt  *time.Time `json:"expiringAt"`
}

type Todo struct {
	TodoID      string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"userId"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Complete    bool      `db:"complete" json:"complete"`
	ExpiringAt  time.Time `db:"expiring_at" json:"expiringAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}
