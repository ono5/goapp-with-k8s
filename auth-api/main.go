// auth-api/main.go
package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// GET用の送信データ
type sendMsg struct {
	Message        string `json:"message"`
	UID            string `json:"uid"`
	Token          string `json:"token"`
	HashedPassword string `json:"hashed_password"`
}

func main() {
	// https://github.com/labstack/echo
	e := echo.New()

	// Routes
	e.GET("/verify-token/:token", tokenHandler)
	e.GET("/token/:hashedPassword/:enteredPassword", hashHandler)
	e.GET("/hashed-password/:password", passwordHandler)

	// Start Server
	e.Logger.Fatal(e.Start(":80"))
}

func tokenHandler(c echo.Context) error {
	token := c.Param("token")
	// ダミートークン
	if token == "abc" {
		return c.JSON(http.StatusOK, sendMsg{Message: "Valid token", UID: "u1"})
	}
	return c.JSON(http.StatusUnauthorized, sendMsg{Message: "token invalid."})
}

func hashHandler(c echo.Context) error {
	hashedPassword := c.Param("hashedPassword")
	enteredPassword := c.Param("enteredPassword")

	// dummy password verification!
	if hashedPassword == enteredPassword+"_hash" {
		const token = "abc"
		return c.JSON(http.StatusOK, sendMsg{Message: "Token created.", Token: token})
	}
	return c.JSON(http.StatusUnauthorized, sendMsg{Message: "Passwords do not match."})
}

func passwordHandler(c echo.Context) error {
	// dummy hashed pw generation!
	enteredPassword := c.Param("password")
	hashedPassword := enteredPassword + "_hash"
	return c.JSON(http.StatusOK, sendMsg{HashedPassword: hashedPassword})
}
