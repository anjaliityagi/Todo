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
