package main

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/steinfletcher/apitest"
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

	t.Run("Normal: Get users list", func(t *testing.T) {
		t.Parallel()

		db, err := newMockDB(t)
		if err != nil {
			t.Error(err)
		}
		apitest.New().
			Handler(newApp(db)).
			Get("/users").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
