package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connet database")
	}

	db.AutoMigrate(&User{})

	userController := NewUserController(db)

	e.GET("/users/:id", userController.getUser)
	e.POST("users", userController.createUser)
	e.Logger.Fatal(e.Start(":8080"))
}
