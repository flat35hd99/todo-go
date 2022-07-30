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

	Name string `json:"name"`
	Age  uint   `json:"age"`
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
	result := controller.db.Take(&user, "id = ?", id)
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

func (controller UserHandler) updateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Structure is not good")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var user User
	result := controller.db.Take(&user, "id = ?", id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if u.Name != "" {
		user.Name = u.Name
	}

	if u.Age != 0 {
		user.Age = u.Age
	}

	result = controller.db.Save(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}

func (controller UserHandler) getUsers(c echo.Context) error {
	var users []User
	result := controller.db.Find(&users)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, struct {
		Users []User `json:"users"`
	}{
		Users: users,
	})
}

func (controller UserHandler) deleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var user User
	result := controller.db.Take(&user, "id = ?", id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	result = controller.db.Delete(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}
