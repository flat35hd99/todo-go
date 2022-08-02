package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
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

	todoHandler := NewTodoHandler(db)
	todoGroup := e.Group("/todos")
	// todoGroup.GET("", todoHandler.getTodos)
	todoGroup.GET("/:id", todoHandler.getTodo)
	todoGroup.POST("", todoHandler.createTodo)

	return e
}

func newDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/tmp/production.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var lambdaAdapter *echoadapter.EchoLambdaV2

func init() {
	log.Print("Cold started")
	db, err := newDB()
	if err != nil {
		panic(err)
	}

	e := newApp(db)
	lambdaAdapter = echoadapter.NewV2(e)
}

func main() {
	log.Print("main started")
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return lambdaAdapter.ProxyWithContext(ctx, req)
	})
}
