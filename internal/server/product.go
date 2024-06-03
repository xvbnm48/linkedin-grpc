package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xvbnm48/linkedin-grpc/internal/dberrors"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"net/http"
)

func (s *EchoServer) GetAllProducts(ctx echo.Context) error {
	vendorID := ctx.QueryParam("vendorId")

	products, err := s.Db.GetAllProducts(ctx.Request().Context(), vendorID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) AddNewProduct(ctx echo.Context) error {
	product := new(models.Product)
	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err.Error())
	}
	fmt.Println("service:`", product)

	newProduct, err := s.Db.AddProduct(ctx.Request().Context(), product)
	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusCreated, newProduct)
}

func (s *EchoServer) GetProductById(ctx echo.Context) error {
	ID := ctx.Param("id")
	product, err := s.Db.GetProductById(ctx.Request().Context(), ID)
	fmt.Println("err", err)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	return ctx.JSON(http.StatusOK, product)
}
