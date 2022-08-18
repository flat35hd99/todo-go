package backend

import (
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

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
			Handler(NewApp(db)).
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
			Handler(NewApp(db)).
			Get("/users").
			Expect(t).
			Status(http.StatusOK).
			End()
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
			Handler(NewApp(db)).
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
			Handler(NewApp(db)).
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

	t.Run("Normal: Get todos", func(t *testing.T) {
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
			Handler(NewApp(db)).
			Get("/todos").
			Expect(t).
			// Status(http.StatusOK).
			Assert(
				jsonpath.Root("$.todos[0]").Equal("title", "test title").Equal("body", "test body").Equal("done", false).Equal("user_id", 1.0).End())
	})

	// User can login by post User.Name and User.Password to /login
	t.Run("Normal: Login", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)

		// Create hashed password using bcrypt
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte("mypassword"), 10)
		if err != nil {
			t.Error(err)
		}

		// Create a user using the hashed password
		type user struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Age            int    `json:"age"`
			Todos          []Todo `json:"todos"`
			HashedPassword string
		}
		u := user{
			Name: "test",
			Age:  30,
			Todos: []Todo{
				{
					Title: "test title",
					Body:  "test body",
					Done:  false,
				},
			},
			HashedPassword: string(hashedPasswordBytes),
		}
		if err := db.Create(&u).Error; err != nil {
			t.Error(err)
		}

		apitest.New().
			Handler(NewApp(db)).
			Post("/auth/login").
			ContentType("application/json").
			Body(fmt.Sprintf(`{"name":"%s","password":"%s"}`, u.Name, "mypassword")).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	// Sign up
	t.Run("Normal: Sign up", func(t *testing.T) {
		t.Parallel()

		db := newMockDB(t)

		// Create a user using the hashed password
		type user struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Age            int    `json:"age"`
			Todos          []Todo `json:"todos"`
			HashedPassword string
		}

		apitest.New().
			Handler(NewApp(db)).
			Post("/auth/signup").
			ContentType("application/json").
			Body(fmt.Sprintf(`{"name":"%s","password":"%s"}`, "myname", "mypassword")).
			Expect(t).
			Status(http.StatusOK).
			End()

		// Expect the user to be created
		u := user{}
		result := db.Where("name = ?", "myname").First(&u)
		if result.Error != nil {
			t.Error(result.Error)
		}
		assert.Equal(t, u.Name, "myname") // Name is set
		// Password is hashed and stored in the database
		// CompareHashAndPassword will return an error if the password is incorrect
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte("mypassword")))
		assert.Error(t, bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte("mypassword2")))
	})
}
