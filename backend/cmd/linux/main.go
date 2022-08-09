package main

import (
	"backend"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

	e := backend.NewApp(db)
	// e.Use(middleware.CORS())
	log.Fatal(e.Start(":8080"))
}
