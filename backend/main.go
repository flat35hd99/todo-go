package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM) // AWS Lambda and GCP Cloud Run uses SIGTERM
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	defer println("Server successfully has been shutdown")
}
