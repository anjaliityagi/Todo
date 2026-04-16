package utils

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

// --------------------
// Error Struct
// --------------------

type Error struct {
	StatusCode    int    `json:"statusCode"`
	Error         string `json:"error"`
	MessageToUser string `json:"messageToUser"`
}

// --------------------
// Response Helpers (Gin)
// --------------------

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

// --------------------
// Request Parsing
// --------------------

func ParseBody(c *gin.Context, out interface{}) error {
	return c.ShouldBindJSON(out)
}

// --------------------
// Password Helpers
// --------------------

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

// --------------------
// Utility Helpers
// --------------------

func ParseBool(str string) bool {
	return str != "" && str != "false"
}

func ParseExpiringAt(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}

	d, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil, err
	}

	if d.Before(time.Now()) {
		return nil, errors.New("invalid time")
	}

	return &d, nil
}
