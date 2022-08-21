package backend

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) AuthHandler {
	return AuthHandler{
		db: db,
	}
}

// Parse the json body and return the user if the password is correct
func (h AuthHandler) Login(c echo.Context) error {
	/* 各種エラーメッセージについて

	## リクエスト形式が間違っているとき
	リクエスト形式が間違っているかどうかは、ユーザーに知られても問題ない。
	本プロダクトはフロントがあるため、解析されたら知られるためである。

	## ユーザーが存在しないとき
	ユーザーが存在するかしないかは、ユーザーに知られてはいけない。
	プロダクトによっては、ユーザーが存在しないことを知られても問題なさそうではあるが、
	本プロダクトにおいてはユーザー一覧のような機能を設ける予定はない。
	need to knowの原則?にあてはまるかはわからないが、ユーザーが存在しないことを伝えても益が少ない。
	そのため、後項のエラーメッセージを使用する。

	## パスワードが間違っているとき
	パスワードが間違っているかどうかは、ユーザーに知られても問題ない。
	パスワードが間違っているかどうかは、ログインできたかで判別可能だからである。
	しかし、前項のユーザーが存在しないときのエラーメッセージと同じコードを用いることで、
	攻撃者がユーザーを間違えたのか、パスワードを間違えたのかを判別できないようにすることができる。

	## 共通化
	上記の3つのエラーメッセージは、共通化することができる。
	いいかえれば、共通化させることで攻撃者に与える情報を少なくすることができる。
	たしかにユーザーとしては、エラーが表示されたとき、ユーザー名を間違えたのかパスワードを間違えたのかを判別できないため、不便である。
	しかし、emailを入力すればパスワードリセットする機能を実装すれば、この不便さはある程度解消できる。
	そのため、Login()の中でのエラーレスポンスはすべてerrorReponseを用いると規定し、
	今後、パスワードリセット機能を実装することとする。
	*/
	errorResponse := echo.NewHTTPError(http.StatusBadRequest, "Authentication failed")

	// Obtain the user credentials from the request body
	type inputLoginInfo struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	input := new(inputLoginInfo)
	if err := c.Bind(input); err != nil {
		return errorResponse
	}

	// Find the user in the database by the given name
	var user User
	result := h.db.Where("name = ?", input.Name).First(&user)
	if result.Error != nil {
		return errorResponse
	}

	// Compare the given password with the stored one
	// If they don't match, return an error
	// Otherwise, return the user
	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password)) != nil {
		return errorResponse
	}
	return c.JSON(http.StatusOK, user)
}

// Sign up a new user
func (h AuthHandler) Signup(c echo.Context) error {

	// Obtain the user credentials from the request body
	type inputSignupInfo struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	input := new(inputSignupInfo)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Check if the user already exists in the database
	var user User
	result := h.db.Where("name = ?", input.Name).First(&user)
	if result.Error == nil {
		return echo.NewHTTPError(http.StatusConflict, "User already exists")
	}

	// If the user does not exist, create a new user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid password")
	}
	user = User{
		Name:           input.Name,
		HashedPassword: string(hashedPassword),
	}
	result = h.db.Create(&user)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}
	return c.JSON(http.StatusOK, user)
}
