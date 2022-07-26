package backend

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	var controller UserHandler
	controller.db = db
	return controller
}

func (controller UserHandler) GetUser(c echo.Context) error {
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

func (handler UserHandler) CreateUser(c echo.Context) error {
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

	result := handler.db.Create(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}

func (handler UserHandler) UpdateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Structure is not good")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var user User
	result := handler.db.Take(&user, "id = ?", id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if u.Name != "" {
		user.Name = u.Name
	}

	if u.Age != 0 {
		user.Age = u.Age
	}

	result = handler.db.Save(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}

func (handler UserHandler) GetUsers(c echo.Context) error {
	var users []User
	result := handler.db.Find(&users)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, struct {
		Users []User `json:"users"`
	}{
		Users: users,
	})
}

func (handler UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var user User
	result := handler.db.Take(&user, "id = ?", id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	result = handler.db.Delete(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}
