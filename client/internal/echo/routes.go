package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Routes() http.Handler {
	e := echo.New()

	v1 := e.Group("v1")
	v1.POST("/sms", s.SmsSendHandler)

	return e
}
