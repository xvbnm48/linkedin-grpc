package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
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
	GetAllProducts(ctx echo.Context) error
	GetAllService(ctx echo.Context) error
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

	pg := s.echo.Group("/products")
	pg.GET("", s.GetAllProducts)

	sg := s.echo.Group("/services")
	sg.GET("", s.GetAllService)
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
