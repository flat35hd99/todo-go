package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

type User struct {
	// Costumized gorm.Model
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index"`

	Name string
	Age  uint
}

func NewUserHandler(db *gorm.DB) UserHandler {
	var controller UserHandler
	controller.db = db
	return controller
}

func (controller UserHandler) getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var user User
	result := controller.db.First(&user, id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (controlloer UserHandler) createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Structure is not good")
	}

	// Use only name and age properties
	// to avoid users to inject illegal value
	user := User{
		Name: u.Name,
		Age:  u.Age,
	}

	result := controlloer.db.Create(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}
