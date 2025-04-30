package main

import (
	"demo_testcontainers/app/users"
	"demo_testcontainers/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	db := db.InitDB()

	e := echo.New()
	userHandler := users.NewUserHandler(db)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users", userHandler.GetUsers)
	e.Logger.Fatal(e.Start(":1323"))
}
