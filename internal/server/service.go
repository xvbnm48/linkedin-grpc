package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *EchoServer) GetAllService(ctx echo.Context) error {
	service, err := s.Db.GetAllService(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, service)
}
