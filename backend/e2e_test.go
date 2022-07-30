package main

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/* Create a mock database
This function will create a mock database and return it.
The database will be closed when the test is finished.
*/
func newMockDB(t *testing.T) (*gorm.DB, error) {
	tmpdir := t.TempDir()

	db, err := gorm.Open(sqlite.Open(filepath.Join(tmpdir, "test.db")), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return db, err
}

func TestE2E(t *testing.T) {
	t.Parallel()

	t.Run("Abnormal: Try to get a users does not exist", func(t *testing.T) {
		t.Parallel()

		db, err := newMockDB(t)
		if err != nil {
			t.Error(err)
		}
		apitest.New().
			Handler(newApp(db)).
			Get("/users/1").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Normal: Get users", func(t *testing.T) {
		t.Parallel()

		db, err := newMockDB(t)
		if err != nil {
			t.Error(err)
		}
		user := User{
			Name: "test",
			Age:  30,
		}
		if err = db.Create(&user).Error; err != nil {
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

		db, err := newMockDB(t)
		if err != nil {
			t.Error(err)
		}
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
}
