package httpinfra

import (
	"context"

	"github.com/labstack/echo"
)

type Server struct {
	port string
	echo *echo.Echo
}

func NewServer(port string) *Server {
	echo := echo.New()
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
