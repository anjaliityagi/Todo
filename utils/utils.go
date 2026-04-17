package utils

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Error struct {
	StatusCode    int    `json:"statusCode"`
	Error         string `json:"error"`
	MessageToUser string `json:"messageToUser"`
}

func RespondJSON(c *gin.Context, statusCode int, body interface{}) {
	c.JSON(statusCode, body)
}

func RespondError(c *gin.Context, statusCode int, err error, messageToUser string) {
	var errString string
	if err != nil {
		errString = err.Error()
	}

	c.JSON(statusCode, Error{
		StatusCode:    statusCode,
		Error:         errString,
		MessageToUser: messageToUser,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
