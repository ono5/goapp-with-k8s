// users-api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

type sendMsg struct {
	Message string `json:"message"`
}

type recvMsg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type hashedPasswordRecvMsg struct {
	HashedPassword string `json:"hashed_password"`
}

var (
	authAddr = os.Getenv("AUTH_ADDRESS") // docker-compose.ymlに設定した環境変数
	authURL  = fmt.Sprintf("http://%s/", authAddr)
)

func main() {
	// https://github.com/labstack/echo
	e := echo.New()

	// Routes
	e.POST("/signup", signupHandler)
	e.POST("/login", loginHandler)

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}

func signupHandler(c echo.Context) error {
	m := new(recvMsg)
	// It's just a dummy service - we don't really care for the email
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Storing the UserInfo failed."})
	}

	if m.Password == "" ||
		len(strings.Trim(m.Password, "")) == 0 ||
		m.Email == "" ||
		len(strings.Trim(m.Email, "")) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, sendMsg{Message: "An email and password needs to be specified!"})
	}

	response, err := http.Get(authURL + "hashed-password/" + m.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Creating the user failed - please try again later."})
	}
	defer response.Body.Close()

	r := new(hashedPasswordRecvMsg)
	if err = json.NewDecoder(response.Body).Decode(r); err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Creating the user failed - please try again later."})
	}

	// since it's a dummy service, we don't really care for the hashed-pw either
	fmt.Println(r.HashedPassword, m.Email)
	return c.JSON(http.StatusCreated, sendMsg{Message: "User created!"})
}

func loginHandler(c echo.Context) error {
	m := new(recvMsg)
	// It's just a dummy service - we don't really care for the email
	if err := c.Bind(m); err != nil {
		log.Fatal(err)
	}

	if m.Password == "" ||
		len(strings.Trim(m.Password, "")) == 0 ||
		m.Email == "" ||
		len(strings.Trim(m.Email, "")) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, sendMsg{Message: "An email and password needs to be specified!"})
	}

	// normally, we'd find a user by email and grab his/ her ID and hashed password
	hashedPassword := m.Password + "_hash"
	response, err := http.Get(authURL + "token/" + hashedPassword + "/" + m.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Internal Server Error"})
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		return c.JSON(http.StatusOK, sendMsg{Message: response.Status})
	}

	return c.JSON(http.StatusUnauthorized, sendMsg{Message: "Logging in failed!"})
}
