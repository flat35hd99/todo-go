package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

/* Create a mock database
This function will create a mock database and return it.
The database will be closed when the test is finished.
*/
func newMockDB(t *testing.T) *gorm.DB {
	tmpdir := t.TempDir()

	db, err := gorm.Open(sqlite.Open(filepath.Join(tmpdir, "test.db")), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	if err = db.AutoMigrate(&User{}, &Todo{}); err != nil {
		t.Error(err)
	}
	return db
}

func TestE2E(t *testing.T) {
	t.Parallel()

	t.Run("Abnormal: Try to get a users does not exist", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)
		apitest.New().
			Handler(newApp(db)).
			Get("/users/1").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Normal: Get users", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)
		user := User{
			Name: "test",
			Age:  30,
		}
		if err := db.Create(&user).Error; err != nil {
			t.Error(err)
		}
		apitest.New().
			Handler(newApp(db)).
			Get("/users").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Normal: Create users", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)
		apitest.New().
			Handler(newApp(db)).
			Post("/users").
			ContentType("application/json").
			Body(`{"name":"test","age":30}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		type user struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		var u user
		assert.NoError(t, db.First(&u).Error)
		assert.Equal(t, u.Name, "test")
		assert.Equal(t, u.Age, 30)

		u_db := User{}
		assert.NoError(t, db.First(&u_db).Error)
		assert.Equal(t, u_db.Name, "test")
		assert.Equal(t, u_db.Age, uint(30))
		assert.NotEmpty(t, u_db.ID)
		assert.NotEmpty(t, u_db.CreatedAt)
		assert.NotEmpty(t, u_db.UpdatedAt)
	})

	t.Run("Normal: Get a todo", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)

		type todo struct {
			Title  string `json:"title"`
			Body   string `json:"body"`
			Done   bool   `json:"done"`
			UserId int    `json:"user_id"`
		}
		type user struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Todos []todo `json:"todos"`
		}
		u := user{
			Name: "test",
			Age:  30,
			Todos: []todo{
				{
					Title: "test title",
					Body:  "test body",
					Done:  false,
				},
			},
		}
		if err := db.Create(&u).Error; err != nil {
			t.Error(err)
		}
		apitest.New().
			Handler(newApp(db)).
			Get("/todos/1").
			Expect(t).
			Status(http.StatusOK).
			Assert(
				jsonpath.Chain().
					Equal("title", "test title").
					Equal("body", "test body").
					Equal("done", false).
					Equal("user_id", 1.0). // TODO: fix this
					End()).
			End()
	})

	t.Run("Normal: Create a todo", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)

		type todo struct {
			Title  string `json:"title"`
			Body   string `json:"body"`
			Done   bool   `json:"done"`
			UserId int    `json:"user_id"`
		}
		type user struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Todos []todo `json:"todos"`
		}
		u := user{
			Name:  "test",
			Age:   30,
			Todos: []todo{},
		}
		if err := db.Create(&u).Error; err != nil {
			t.Error(err)
		}
		apitest.New().
			Handler(newApp(db)).
			Post("/todos").
			ContentType("application/json").
			Body(fmt.Sprintf(`{"title":"test title","body":"test body","done":false,"user_id":%d}`, u.ID)).
			Expect(t).
			Status(http.StatusOK).
			End()

		td := Todo{}
		assert.NoError(t, db.First(&td).Error)
		assert.Equal(t, td.Title, "test title")
		assert.Equal(t, td.Body, "test body")
		assert.Equal(t, td.Done, false)
		assert.Equal(t, td.UserID, u.ID)
	})
}
