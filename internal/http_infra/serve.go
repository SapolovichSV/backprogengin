package httpinfra

import (
	"context"

	_ "github.com/SapolovichSV/backprogeng/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	port string
	echo *echo.Echo
}

func NewServer(port string) *Server {
	echo := echo.New()
	echo.GET("/swagger/*", echoSwagger.WrapHandler)
	echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	}))
	echo.Use(middleware.Logger())
	return &Server{
		port: port,
		echo: echo,
	}
}
func (s *Server) GetRouter() *echo.Router {
	return s.echo.Router()
}
func (s *Server) Start() error {
	return s.echo.Start(":" + s.port)
}
func (s *Server) Stop(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
