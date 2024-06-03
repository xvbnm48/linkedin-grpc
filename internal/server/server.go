package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xvbnm48/linkedin-grpc/internal/models"
	"log"
	"net/http"

	"github.com/xvbnm48/linkedin-grpc/internal/database"
)

type Server interface {
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error

	GetAllCustomer(ctx echo.Context) error
	AddNewCustomer(ctx echo.Context) error
	GetCustomerById(ctx echo.Context) error

	GetAllProducts(ctx echo.Context) error
	AddNewProduct(ctx echo.Context) error
	GetProductById(ctx echo.Context) error

	GetAllService(ctx echo.Context) error
	AddNewService(ctx echo.Context) error
	GetServiceById(ctx echo.Context) error

	GetAllVendor(ctx echo.Context) error
	AddNewVendor(ctx echo.Context) error
	GetVendorById(ctx echo.Context) error
}

type EchoServer struct {
	echo *echo.Echo
	Db   database.DatabaseClient
}

func NewEchoServer(db database.DatabaseClient) Server {
	server := &EchoServer{
		echo: echo.New(),
		Db:   db,
	}

	server.registerRoutes()
	return server

}

func (s *EchoServer) Start() error {
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	if err := s.echo.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server shutdown with error: %v", err)
	}
	return nil
}
func (s *EchoServer) registerRoutes() {
	s.echo.GET("/readiness", s.Readiness)
	s.echo.GET("/liveness", s.Liveness)
	cg := s.echo.Group("/customers")
	cg.GET("", s.GetAllCustomer)
	cg.POST("", s.AddNewCustomer)
	cg.GET("/:id", s.GetCustomerById)

	pg := s.echo.Group("/products")
	pg.GET("", s.GetAllProducts)
	pg.POST("", s.AddNewProduct)
	pg.GET("/:id", s.GetProductById)

	sg := s.echo.Group("/services")
	sg.GET("", s.GetAllService)
	sg.POST("", s.AddNewService)
	sg.GET("/:id", s.GetServiceById)

	vg := s.echo.Group("/vendors")
	vg.GET("", s.GetAllVendor)
	vg.POST("", s.AddNewVendor)
	vg.GET("/:id", s.GetVendorById)
}
func (s *EchoServer) Readiness(ctx echo.Context) error {
	ready := s.Db.Ready()
	fmt.Println("ready:", ready)
	if ready {
		return ctx.JSON(http.StatusOK, models.Health{
			Status: "OK",
		})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{
		Status: "Failure",
	})
}

func (s *EchoServer) Liveness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{
		Status: "OK",
	})
}
