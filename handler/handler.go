package handler

import (
	"Todo-Server/database/dbHelper"
	"Todo-Server/models"
	"Todo-Server/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todoRequest models.CreateTodo

	userID := "550e8400-e29b-41d4-a716-446655440000"

	if err := utils.ParseBody(c, &todoRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	if err := utils.Validate.Struct(todoRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "validation failed")
		return
	}

	if todoRequest.ExpiringAt.Before(time.Now()) {
		utils.RespondError(c, http.StatusBadRequest, nil, "provided time and date is wrong")
		return
	}

	todo, err := dbHelper.CreateTodo(
		userID,
		todoRequest.Name,
		todoRequest.Description,
		todoRequest.ExpiringAt,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, todo)
}

func GetAllTodos(c *gin.Context) {
	fmt.Println("hello")
	completeStr := c.Query("status")
	expiringAtStr := c.Query("expiringAt")
	search := c.Query("search")

	userID := "550e8400-e29b-41d4-a716-446655440000"

	if expiringAtStr != "" {
		d, err := time.Parse("2006-01-02", expiringAtStr)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, err, "invalid date")
			return
		}
		if d.Before(time.Now()) {
			expiringAtStr = ""
		}
	}
	fmt.Println("hello")
	todos, err := dbHelper.GetTodos(userID, search, expiringAtStr, completeStr)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "Failed to fetch todos")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"todos": todos,
	})
}

func GetTodoById(c *gin.Context) {
	todoID := c.Param("id")

	if todoID == "" {
		utils.RespondError(c, http.StatusBadRequest, nil, "todo id is required")
		return
	}

	userID := "550e8400-e29b-41d4-a716-446655440000"

	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, todo)
}

func DeleteTodoById(c *gin.Context) {
	todoID := c.Param("id")

	if todoID == "" {
		utils.RespondError(c, http.StatusBadRequest, nil, "required todo ID")
		return
	}

	userID := "550e8400-e29b-41d4-a716-446655440000"

	err := dbHelper.DeleteTodoById(userID, todoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondError(c, http.StatusNotFound, err, "todo not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to delete todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"message": "todo deleted successfully",
	})
}

func UpdateTodoById(c *gin.Context) {
	todoID := c.Param("id")

	if todoID == "" {
		utils.RespondError(c, http.StatusBadRequest, nil, "required todo id")
		return
	}

	var todo models.UpdateTodoRequest

	if err := utils.ParseBody(c, &todo); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	if err := utils.Validate.Struct(todo); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "validation failed")
		return
	}

	userID := "550e8400-e29b-41d4-a716-446655440000"
	err := dbHelper.UpdateTodoById(
		todo.Name,
		todo.Description,
		todo.Complete,
		todo.ExpiringAt,
		todoID,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondError(c, http.StatusNotFound, err, "todo not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to update todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"message": "updated successfully",
	})
}
