package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewTodoHandler(db *gorm.DB) TodoHandler {
	var controller TodoHandler
	controller.db = db
	return controller
}

type TodoHandler struct {
	db *gorm.DB
}

func (handler TodoHandler) getTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Irregal ID")
	}

	var todo Todo
	result := handler.db.Take(&todo, "id = ?", id)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	return c.JSON(http.StatusOK, todo)
}

func (handler TodoHandler) createTodo(c echo.Context) error {
	inputTodo := new(Todo)
	if err := c.Bind(&inputTodo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	todo := Todo{
		Title:  inputTodo.Title,
		Body:   inputTodo.Body,
		Done:   inputTodo.Done,
		UserID: inputTodo.UserID,
	}

	result := handler.db.Create(&todo)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create todo")
	}

	return c.JSON(http.StatusOK, todo)
}
