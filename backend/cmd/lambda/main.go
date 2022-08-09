package main

import (
	"backend"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

	e := backend.NewApp(db)
	lambdaAdapter = echoadapter.NewV2(e)
}

func main() {
	log.Print("main started")
	lambda.Start(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return lambdaAdapter.ProxyWithContext(ctx, req)
	})
}
