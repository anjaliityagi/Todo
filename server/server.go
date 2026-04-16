package server

import (
	"Todo-Server/handler"
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
		v1.POST("/todos", handler.CreateTodo)
		v1.GET("/todos", handler.GetAllTodos)
		v1.GET("/todos/:id", handler.GetTodoById)
		v1.PUT("/todos/:id", handler.UpdateTodoById)
		v1.DELETE("/todos/:id", handler.DeleteTodoById)
		v1.POST("/logout", handler.Logout)
	}

	return &Server{
		Router: router,
	}
}
