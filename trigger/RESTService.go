package trigger

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	conn *echo.Echo
}

func StartServer() *Server {
	e := echo.New()

	e.Start(":55555")

	e.POST("/alert", alertController)

	return &Server{e}
}

func alertController(c echo.Context) error {
	msg := new(DetectedStatus)
	err := c.Bind(&msg)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.JSON(http.StatusOK, "Success")
}
