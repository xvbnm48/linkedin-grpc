package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *EchoServer) GetAllCustomer(ctx echo.Context) error {
	emailAddress := ctx.QueryParam("email")

	customers, err := s.Db.GetAllCustomers(ctx.Request().Context(), emailAddress)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, customers)
}
