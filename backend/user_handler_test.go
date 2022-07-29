package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connet database")
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}

	userJSON := `{"name": "Bob", "age": 22}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	h := NewUserHandler(db)
	if assert.NoError(t, h.createUser(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var u struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Age       int    `json:"age"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"Updated_at"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
			t.Error(err)
		}
		assert.Equal(t, "Bob", u.Name)
		assert.Equal(t, 22, u.Age)
		assert.NotEmpty(t, u.ID)
		assert.NotEmpty(t, u.CreatedAt)
		assert.NotEmpty(t, u.UpdatedAt)
	}
}
