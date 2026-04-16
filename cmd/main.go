package main

import (
	"Todo-Server/database"
	"Todo-Server/server"
	"fmt"
	"log"
	"os"
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

	database.ConnectandMigrate(
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		database.SSLMode(sslMode),
	)

	fmt.Println("Database connected")

	srv := server.SetupRoutes()

	fmt.Println("Server running at http://localhost:8080")

	if err := srv.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
