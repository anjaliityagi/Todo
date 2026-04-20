package handler

import (
	"Todo-Server/database/dbHelper"
	"Todo-Server/models"
	"Todo-Server/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todoRequest models.CreateTodo

	token := c.GetHeader("Authorization")

	userID, err := dbHelper.GetUserIDBySession(token)
	fmt.Println(userID)
	if err != nil {
		utils.RespondError(c, 401, err, "Token Error")
		return
	}

	if err := c.ShouldBindJSON(&todoRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	if todoRequest.ExpiringAt.Before(time.Now()) {
		utils.RespondError(c, http.StatusBadRequest, nil, "provided time and date is wrong")
		return
	}

	err = dbHelper.CreateTodo(
		userID,
		todoRequest.Name,
		todoRequest.Description,
		todoRequest.ExpiringAt,
	)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(c, http.StatusCreated, "todo created successfully")
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
		utils.RespondError(c, http.StatusBadRequest, nil, "user already exists")
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
		//"token":   sessionID,
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
		utils.RespondError(c, http.StatusInternalServerError, err, "wrong password")
		return
	}

	sessionID, err := dbHelper.CreateUserSession(userDetail.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "failed to create user session")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"token": sessionID,
	})
}

func Logout(c *gin.Context) {

	token := c.GetHeader("Authorization")

	if token == "" {
		utils.RespondError(c, http.StatusUnauthorized, nil, "missing token")
		return
	}

	err := dbHelper.DeleteSessionByToken(token)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err, "invalid session")
		return
	}

	utils.RespondJSON(c, http.StatusOK, gin.H{
		"message": " user logged out successfully",
	})
}

func UpdateTodo(c *gin.Context) {
	var updateTodo models.UpdateTodo
	if err := c.ShouldBindJSON(&updateTodo); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	token := c.GetHeader("Authorization")
	userID, err := dbHelper.GetUserIDBySession(token)
	todoID := c.Params.ByName("todo-id")
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err, "invalid session")
	}
	err = dbHelper.UpdateTodo(todoID, userID, updateTodo.Name, updateTodo.Description, updateTodo.ExpiringAt, updateTodo.Complete)
	if err != nil {
		utils.RespondError(c, 500, err, "invalid request")
	}

	utils.RespondJSON(c, http.StatusCreated, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodo(c *gin.Context) {
	token := c.GetHeader("Authorization")
	userID, err := dbHelper.GetUserIDBySession(token)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err, "invalid session")
	}
	todoID := c.Params.ByName("todo-id")
	err = dbHelper.DeleteTodo(todoID, userID)
	if err != nil {
		utils.RespondError(c, 500, err, "invalid request")
	}
	utils.RespondJSON(c, http.StatusCreated, gin.H{
		"message": "todo deleted successfully",
	})
}

func FetchTodoById(c *gin.Context) {

	token := c.GetHeader("Authorization")
	userID, err := dbHelper.GetUserIDBySession(token)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err, "invalid session")
		return
	}
	id := c.Params.ByName("todo-id")
	var todo models.Todo
	todo, err = dbHelper.GetTodoById(id, userID)

	if err != nil {
		utils.RespondError(c, 500, err, "invalid request")
		return
	}
	utils.RespondJSON(c, http.StatusOK, todo)

}

func GetAllTodos(c *gin.Context) {
 
	search := c.Query("search")
	status := c.Query("status")
	date := c.Query("expiringAt")

	token := c.GetHeader("Authorization")

	userID, err := dbHelper.GetUserIDBySession(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid session"})
		return
	}

	todos, err := dbHelper.GetTodos(userID, search, date, status)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.RespondJSON(c, http.StatusOK, todos)
}
