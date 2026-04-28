package handler

import (
	"Todo-Server/database/dbHelper"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FetchAllUsers(c *gin.Context) {

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	total, err := dbHelper.FetchUsersCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch user count"})
		return
	}

	users, err := dbHelper.FetchAllUsers(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{
		"total": total,
		"users": users,
	})
}

func GetAllTodos(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	total, err := dbHelper.FetchTodosCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todo count"})
		return
	}

	todos, err := dbHelper.FetchAllTodos(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(200, gin.H{
		"total": total,
		"todos": todos,
	})
}

func SuspendUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	if userID == currentUserID {
		c.JSON(400, gin.H{"error": "cannot suspend yourself"})
		return
	}

	err := dbHelper.SuspendUser(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to suspend user"})
		return
	}

	c.JSON(200, gin.H{"message": "user suspended successfully"})
}
