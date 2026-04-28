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

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())

		userRoutes := auth.Group("/")
		userRoutes.Use(middleware.RequireRole("user", "admin"))

		userRoutes.POST("/todo", handler.CreateTodo)
		userRoutes.PUT("/updatetodo/:todo-id", handler.UpdateTodoById)
		userRoutes.DELETE("/deletetodo/:todo-id", handler.DeleteTodo)
		userRoutes.GET("/todo/:todo-id", handler.FetchTodoById)
		userRoutes.GET("/todos", handler.FetchAllTodos)

		adminRoutes := auth.Group("/admin")
		adminRoutes.Use(middleware.RequireRole("admin"))

		adminRoutes.GET("/users", handler.FetchAllUsers)
		adminRoutes.GET("/todos", handler.FetchAllTodosForAllUsers)
		adminRoutes.PUT("/users/:user-id/toggleSuspend", handler.ToggleSuspendTx)

		auth.PUT("/logout", handler.Logout)
	}

	return &Server{
		Router: router,
	}
}
