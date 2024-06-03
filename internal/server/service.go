package server

import (
	"github.com/labstack/echo/v4"
	"github.com/xvbnm48/linkedin-grpc/internal/dberrors"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"net/http"
)

func (s *EchoServer) GetAllService(ctx echo.Context) error {
	service, err := s.Db.GetAllService(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, service)
}

func (s *EchoServer) AddNewService(ctx echo.Context) error {
	service := new(models.Service)
	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err.Error())
	}

	newService, err := s.Db.AddService(ctx.Request().Context(), service)
	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusCreated, newService)
}

func (s *EchoServer) GetServiceById(ctx echo.Context) error {
	ID := ctx.Param("id")
	service, err := s.Db.GetServiceByID(ctx.Request().Context(), ID)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

	}

	return ctx.JSON(http.StatusOK, service)
}
