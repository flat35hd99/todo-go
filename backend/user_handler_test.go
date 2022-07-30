package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	e := echo.New()

	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "test.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
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
			UpdatedAt string `json:"updated_at"`
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

func seedUser(db *gorm.DB) {
	for _, u := range []User{
		{Name: "Bob", Age: 22},
		{Name: "Alice", Age: 33},
		{Name: "山田", Age: 44},
	} {
		db.Create(&u)
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	e := echo.New()

	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "test.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}
	seedUser(db)

	t.Run("Normal: Get a user", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")
		h := NewUserHandler(db)
		if assert.NoError(t, h.getUser(ctx)) {
			assert.Equal(t, 200, rec.Code)
			var u struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Age       int    `json:"age"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
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
	})

	t.Run("Abnormal: Try to get a users does not exist", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("0")
		h := NewUserHandler(db)
		err := h.getUser(ctx)
		assert.Error(t, err)
	})

	t.Run("Abnormal: Try to inject SQL", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("\\' OR 1=1 --")
		h := NewUserHandler(db)
		err := h.getUser(ctx)
		assert.Error(t, err)
	})

	t.Run("Abnormal: Wrong id format", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("a")
		h := NewUserHandler(db)
		err := h.getUser(ctx)
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	e := echo.New()

	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	db, err := gorm.Open(sqlite.Open(filepath.Join(dir, "test.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}
	seedUser(db)

	t.Run("Normal: Update name of a user", func(t *testing.T) {
		t.Parallel()
		userJSON := `{"name": "Captain", "age": 22}`
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")
		h := NewUserHandler(db)
		if assert.NoError(t, h.updateUser(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var u struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Age       int    `json:"age"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
				t.Error(err)
			}
			assert.Equal(t, "Captain", u.Name)
			assert.Equal(t, 22, u.Age)
			assert.NotEmpty(t, u.ID)
			assert.NotEmpty(t, u.CreatedAt)
			assert.NotEmpty(t, u.UpdatedAt)

		}
	})
}
