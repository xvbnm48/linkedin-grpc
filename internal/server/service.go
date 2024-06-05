package server

import (
	"fmt"
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
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}

	}

	return ctx.JSON(http.StatusOK, service)
}

func (s *EchoServer) UpdateService(ctx echo.Context) error {
	ID := ctx.Param("id")
	service := new(models.Service)
	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	fmt.Println("this is a body", service)

	if ID != service.ServiceID {
		return ctx.JSON(http.StatusBadRequest, "ID in the path does not match the ID in the body")
	}
	service, err := s.Db.UpdateService(ctx.Request().Context(), service)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, service)
}

func (s *EchoServer) DeleteService(ctx echo.Context) error {
	id := ctx.Param("id")

	err := s.Db.DeleteService(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Service deleted successfully")
}
