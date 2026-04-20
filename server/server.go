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
		v1.PUT("/updatetodo/:todo-id", handler.UpdateTodo)
		v1.DELETE("/deletetodo/:todo-id", handler.DeleteTodo)
		v1.PUT("/logout", handler.Logout)
		v1.GET("/todo/:todo-id", handler.FetchTodoById)
		v1.GET("/todos", handler.GetAllTodos)

	}

	return &Server{
		Router: router,
	}
}
