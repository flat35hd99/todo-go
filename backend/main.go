package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newApp(db *gorm.DB) *echo.Echo {
	e := echo.New()

	userHandler := NewUserHandler(db)
	userGroup := e.Group("/users")
	userGroup.GET("", userHandler.getUsers)
	userGroup.GET("/:id", userHandler.getUser)
	userGroup.POST("", userHandler.createUser)
	userGroup.PATCH("/:id", userHandler.updateUser)
	userGroup.DELETE("/:id", userHandler.deleteUser)

	return e
}

func newDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("production.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := newDB()
	if err != nil {
		panic(err)
	}
	e := newApp(db)
	e.Logger.Fatal(e.Start(":8080"))
}
