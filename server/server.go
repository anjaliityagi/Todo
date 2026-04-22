package server

import (
	"Todo-Server/handler"
	"Todo-Server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func SetupRoutes() *Server {

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "Server is Running",
			})
		})
		v1.POST("/register", handler.RegisterUser)
		v1.POST("/login", handler.LoginUser)

		protected := v1.Group("/")

		protected.Use(middleware.AuthMiddleware())
		todoPath := protected.Group("/todo")

		todoPath.POST("", handler.CreateTodo)
		todoPath.PUT("/:todo-id", handler.UpdateTodoById)
		todoPath.DELETE("/:todo-id", handler.DeleteTodo)
		todoPath.GET("/:todo-id", handler.FetchTodoById)

		protected.GET("/todos", handler.FetchAllTodos)
		protected.PUT("/logout", handler.Logout)

	}

	return &Server{
		Router: router,
	}
}
