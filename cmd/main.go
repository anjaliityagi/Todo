package main

import (
	"Todo-Server/database"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func main() {

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5434")
	dbUser := getEnv("DB_USER", "local")
	dbPassword := getEnv("DB_PASSWORD", "local")
	dbName := getEnv("DB_NAME", "mercury-dev")
	sslMode := getEnv("DB_SSLMODE", string(database.SSLModeDisable))
	//serverPort := getEnv("SERVER_PORT", "8080")

	database.ConnectandMigrate(
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		database.SSLMode(sslMode),
	)
	router := gin.Default()
	fmt.Println("server is running")
	router.GET("/hello", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	router.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "health check",
		})
	})
	fmt.Println("Server running at http://localhost:8080/hello")
	fmt.Println("server started at:8080", 8080)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
