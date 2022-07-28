package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name string
	Age  uint
}

func NewUserController(db *gorm.DB) UserController {
	var controller UserController
	controller.db = db
	return controller
}

func (controller UserController) getUser(c echo.Context) error {
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

func (controlloer UserController) createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Structure is not good")
	}

	user := User{
		Name: u.Name,
		Age:  u.Age,
	}

	result := controlloer.db.Create(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User could not be created")
	}

	return c.JSON(http.StatusOK, user)
}
