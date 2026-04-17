package handler

import (
	"Todo-Server/database/dbHelper"
	"Todo-Server/models"
	"Todo-Server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var todoRequest models.CreateTodo
	userID := c.Param("user-id")

	if err := c.ShouldBindJSON(&todoRequest); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	if todoRequest.ExpiringAt.Before(time.Now()) {
		utils.RespondError(c, http.StatusBadRequest, nil, "provided time and date is wrong")
		return
	}

	err := dbHelper.CreateTodo(
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
