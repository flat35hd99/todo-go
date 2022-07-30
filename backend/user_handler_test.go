package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	e := echo.New()

	db, err := newMockDB(t)
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

/* Insert users and get them back
{ID: 1, Name: "Bob", Age: 22},
{ID: 2, Name: "Alice", Age: 33},
{ID: 3, Name: "山田", Age: 44},
*/
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

	db, err := newMockDB(t)
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

	db, err := newMockDB(t)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}
	seedUser(db)

	t.Run("Normal: Update name of a user", func(t *testing.T) {
		t.Parallel()

		// Update bob's name to Captain
		userJSON := `{"name": "Captain"}`
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

	t.Run("Normal: Update Age of a user", func(t *testing.T) {
		t.Parallel()

		// Update alice's age from 33 to 30
		userJSON := `{"age": 30}`
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("2")
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
			assert.Equal(t, "Alice", u.Name)
			assert.Equal(t, 30, u.Age)
			assert.NotEmpty(t, u.ID)
			assert.NotEmpty(t, u.CreatedAt)
			assert.NotEmpty(t, u.UpdatedAt)

		}
	})
}

func TestGetUsers(t *testing.T) {
	t.Parallel()
	e := echo.New()

	db, err := newMockDB(t)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}
	seedUser(db)

	t.Run("Normal: Get all users", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users")
		h := NewUserHandler(db)
		if assert.NoError(t, h.getUsers(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			var data struct {
				Users []struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					Age       int    `json:"age"`
					CreatedAt string `json:"created_at"`
					UpdatedAt string `json:"updated_at"`
				} `json:"users"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
				t.Error(err)
			}
			users := data.Users
			assert.Equal(t, 3, len(users))
			assert.Equal(t, "Bob", users[0].Name)
			assert.Equal(t, "Alice", users[1].Name)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	e := echo.New()

	db, err := newMockDB(t)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		t.Errorf("%v\n", err)
	}
	seedUser(db)

	t.Run("Normal: Delete a user", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)

		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetPath("/users/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")
		h := NewUserHandler(db)
		if assert.NoError(t, h.deleteUser(ctx)) {
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

		// Check if the user is deleted from the database
		var users []User
		db.Find(&users)
		assert.Equal(t, 2, len(users))
		found := false
		for _, u := range users {
			if u.ID == 1 {
				found = true
			}
		}
		assert.False(t, found)
	})
}
