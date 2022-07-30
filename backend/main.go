package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connet database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	userHandler := NewUserHandler(db)
	userGroup := e.Group("/users")
	userGroup.GET("", userHandler.getUsers)
	userGroup.GET("/:id", userHandler.getUser)
	userGroup.POST("", userHandler.createUser)
	userGroup.PATCH("/:id", userHandler.updateUser)
	userGroup.DELETE("/:id", userHandler.deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
