package main

import (
	"backend"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
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
	db, err := gorm.Open(sqlite.Open("/tmp/production.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&backend.User{}, &backend.Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var lambdaAdapter *echoadapter.EchoLambdaV2

func init() {
	log.Print("Cold started")
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	e := NewApp(db)
	lambdaAdapter = echoadapter.NewV2(e)
}

func main() {
	log.Print("main started")
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return lambdaAdapter.ProxyWithContext(ctx, req)
	})
}
