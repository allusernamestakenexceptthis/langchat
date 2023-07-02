package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HomeHandler is a handler for the home page
func login(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Worlda!")
}

type User struct {
	Username string
}
