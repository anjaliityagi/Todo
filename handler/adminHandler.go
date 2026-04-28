package handler

import (
	"Todo-Server/database/dbHelper"
	"Todo-Server/utils"
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
		utils.RespondError(c, 500, err, "failed to fetch users")
		return
	}

	c.JSON(200, gin.H{
		"total": total,
		"users": users,
	})
}

func FetchAllTodosForAllUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	total, err := dbHelper.FetchTodosCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todo count"})
		return
	}

	todos, err := dbHelper.FetchAllTodosForAllUsers(limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch todos"})
		return
	}

	c.JSON(200, gin.H{
		"total": total,
		"todos": todos,
	})
}

func ToggleSuspendTx(c *gin.Context) {
	userID := c.Param("user-id")
	currentUserID := c.GetString("userID")
	suspendStr := c.Query("isSuspend")
	if suspendStr == "" {
		c.JSON(400, gin.H{"error": "Empty suspend"})
		return
	}
	suspend, err := strconv.ParseBool(suspendStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid boolean value"})
		return
	}
	//suspend := c.Query("suspend")

	if userID == currentUserID {
		c.JSON(400, gin.H{"error": "cannot suspend/unsuspend yourself"})
		return
	}

	err = dbHelper.SuspendUserTx(userID, suspend)
	if err != nil {
		utils.RespondError(c, 500, err, "failed to suspend user")
		return
	}

	c.JSON(200, gin.H{"message": "user suspended/unsuspended"})
}
