package handler

import (
	"Todo-Server/database/dbHelper"
	"Todo-Server/models"
	"Todo-Server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todoRequest models.CreateTodo

	userID := c.GetString("userID")

	if err := c.ShouldBindJSON(&todoRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	if todoRequest.ExpiringAt.Before(time.Now()) {
		utils.RespondError(c, http.StatusBadRequest, nil, "provided time and date is wrong")
		return
	}

	todoId, err := dbHelper.CreateTodo(
		userID,
		todoRequest.Name,
		todoRequest.Description,
		todoRequest.ExpiringAt,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, gin.H{
		"TodoID":  todoId,
		"message": "todo created successfully",
	})
}

func RegisterUser(c *gin.Context) {
	var registerUser models.RegisterUser

	if err := c.ShouldBindJSON(&registerUser); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	userExist, err := dbHelper.IsUserExists(registerUser.Email)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to check user existence")
		return
	}
	if userExist {
		utils.RespondError(c, http.StatusConflict, nil, "user already exists")
		return
	}

	hashPassword, err := utils.HashPassword(registerUser.Password)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed while hashing password")
		return
	}

	err = dbHelper.CreateUser(registerUser.Name, registerUser.Email, hashPassword)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create user")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

func LoginUser(c *gin.Context) {
	var req models.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "invalid request body")
		return
	}

	userDetail, err := dbHelper.GetUserByEmail(req.Email)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err, "invalid credentials")
		return
	}

	err = utils.CheckPassword(userDetail.Password, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusForbidden, err, "wrong password")
		return
	}

	sessionID, err := dbHelper.CreateUserSession(userDetail.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create user session")
		return
	}

	//utils.RespondJSON(c, http.StatusOK, gin.H{
	//	"token": sessionID,
	//})

	token, err := utils.GenerateJWT(userDetail.ID, userDetail.Role, sessionID)
	if err != nil {
		utils.RespondError(c, 500, err, "could not generate token")
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func Logout(c *gin.Context) {

	token := c.GetString("sessionId")

	if token == "" {
		utils.RespondError(c, http.StatusUnauthorized, nil, "missing token")
		return
	}

	err := dbHelper.DeleteSessionByToken(token)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "invalid session")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"message": "user logged out successfully",
	})
}

func UpdateTodoById(c *gin.Context) {
	var updateTodo models.UpdateTodo

	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	userID := c.GetString("userID")
	todoID := c.Param("todo-id")

	err := dbHelper.UpdateTodoByID(
		todoID,
		userID,
		updateTodo)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to update todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{

		"message": "todo updated successfully",
	})
}

func DeleteTodo(c *gin.Context) {

	userID := c.GetString("userID")

	todoID := c.Param("todo-id")

	err := dbHelper.DeleteTodo(todoID, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to delete todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{

		"message": "todo deleted successfully",
	})
}

func FetchTodoById(c *gin.Context) {

	userID := c.GetString("userID")

	id := c.Param("todo-id")

	todo, err := dbHelper.FetchTodoById(id, userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to fetch todo")
		return
	}

	utils.RespondJSON(c, http.StatusOK, todo)
}

func FetchAllTodos(c *gin.Context) {

	search := c.Query("search")
	status := c.Query("status")
	date := c.Query("expiringAt")

	limitStr := c.Query("limit")
	pageStr := c.Query("page")

	limit := 1
	page := 1

	if limitStr != "" {
		limitTemp, err := strconv.Atoi(limitStr)
		if err == nil && limitTemp > 0 {
			limit = limitTemp
		}
	}

	if pageStr != "" {
		pageTemp, err := strconv.Atoi(pageStr)
		if err == nil && pageTemp > 0 {
			page = pageTemp
		}
	}

	userID := c.GetString("userID")

	total, err := dbHelper.FetchTotalTodoCount(userID)

	if err != nil {
		utils.RespondError(c, 500, err, "couldn't fetch total count")
		return
	}
	offset := (page - 1) * limit
	totalPages := (total + limit - 1) / limit

	todos, err := dbHelper.FetchTodos(userID, search, date, status, limit, page, offset)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to fetch todos")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"allTodos":   todos,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages})
}
