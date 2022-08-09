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
	userGroup.POST("", userHandler.CreateUser)
	userGroup.PATCH("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)

	todoHandler := NewTodoHandler(db)
	todoGroup := e.Group("/todos")
	todoGroup.GET("", todoHandler.GetTodos)
	todoGroup.GET("/:id", todoHandler.GetTodo)
	todoGroup.POST("", todoHandler.CreateTodo)

	return e
}
