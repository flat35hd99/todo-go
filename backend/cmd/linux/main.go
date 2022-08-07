package main

import (
	"backend"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *echo.Echo {
	e := echo.New()

	userHandler := backend.NewUserHandler(db)
	userGroup := e.Group("/users")
	userGroup.GET("", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUser)
	userGroup.POST("", userHandler.CreateUser)
	userGroup.PATCH("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)

	todoHandler := backend.NewTodoHandler(db)
	todoGroup := e.Group("/todos")
	// todoGroup.GET("", todoHandler.getTodos)
	todoGroup.GET("/:id", todoHandler.GetTodo)
	todoGroup.POST("", todoHandler.CreateTodo)

	return e
}

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&backend.User{}, &backend.Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	e := NewApp(db)
	log.Fatal(e.Start(":8080"))
}
