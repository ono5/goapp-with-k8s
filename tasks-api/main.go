// tasks-api/main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	authAddr = os.Getenv("AUTH_ADDRESS") // docker-compose.ymlに設定した環境変数
	authURL  = fmt.Sprintf("http://%s/", authAddr)
	filePath = os.Getenv("TASKS_FOLDER")
)

// GET用の送信データ
type sendMsg struct {
	Message string `json:"message"`
	Tasks   []task `json:"tasks"`
}

type uidRecvMsg struct {
	UID string `json:"uid"`
}

type task struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func extractAndVerifyToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(fmt.Sprintf("%s", "No token provided."))
	}

	// expects Bearer TOKEN
	token := strings.Split(authHeader, " ")[1]

	response, err := http.Get(authURL + "verify-token/" + token)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	r := new(uidRecvMsg)
	if err = json.NewDecoder(response.Body).Decode(r); err != nil {
		return "", err
	}

	return r.UID, nil
}

func main() {
	// https://github.com/labstack/echo
	e := echo.New()

	// CORSの設定
	e.Use(middleware.CORS())

	// Routes
	e.GET("/tasks", getTaskHandler)
	e.POST("/tasks", postTaskHandler)

	// Start Server
	e.Logger.Fatal(e.Start(":8000"))
}

func getTaskHandler(c echo.Context) error {
	uid, err := extractAndVerifyToken(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, sendMsg{Message: "Could not verify token."})
	}
	fmt.Println(uid) // このuidは特に使わない
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Loading the tasks failed."})
	}
	strData := string(bytes)
	entries := strings.Split(strData, "TASK_SPLIT")
	fmt.Println(entries)
	var tasks []task
	for _, entry := range entries {
		fmt.Println(entry)
		if entry != "" {
			t := new(task)
			err := json.Unmarshal([]byte(entry), &t)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Faild JSON Parse."})
			}
			tasks = append(tasks, *t)
		}
	}
	return c.JSON(http.StatusOK, sendMsg{Message: "Tasks loaded.", Tasks: tasks})
}

func postTaskHandler(c echo.Context) error {
	uid, err := extractAndVerifyToken(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, sendMsg{Message: "Could not verify token."})
	}
	fmt.Println(uid) // このuidは特に使わない
	t := new(task)
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Storing the task failed."})
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Storing the text failed."})
	}
	defer f.Close()
	json, err := json.Marshal(t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, sendMsg{Message: "Faild JSON Parse."})
	}
	f.WriteString(string(json) + "TASK_SPLIT")

	return c.JSON(http.StatusCreated, t)
}
