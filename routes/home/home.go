package home

import (
	"github.com/labstack/echo/v4"
)

// HomeHandler is a handler for the home page
func Home(c echo.Context) error {
	//serve file front/index
	return c.File("front/index.html")

}
