package backend

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *echo.Echo {
	e := echo.New()

	userHandler := NewUserHandler(db)
	userGroup := e.Group("/users")
	userGroup.GET("", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUser)
	// userGroup.POST("", userHandler.CreateUser) user should be created by /auth/signup
	userGroup.PATCH("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)

	todoHandler := NewTodoHandler(db)
	todoGroup := e.Group("/todos")
	todoGroup.GET("", todoHandler.GetTodos)
	todoGroup.GET("/:id", todoHandler.GetTodo)
	todoGroup.POST("", todoHandler.CreateTodo)

	authHandler := NewAuthHandler(db)
	authGroup := e.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/signup", authHandler.Signup)

	return e
}
