package httpserv

import (
	"context"
	"docs-hub/internal/cloud"

	"docs-hub/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "docs-hub/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type ServerHttp struct {
	config *server.Config
	cloud  *cloud.DocumentHub
	server *echo.Echo
}

func Init(conf *server.Config, cloud *cloud.DocumentHub) *server.Server {
	httpServer := &ServerHttp{
		config: conf,
		cloud:  cloud,
		server: echo.New(),
	}

	return &server.Server{Server: httpServer}
}

func (s *ServerHttp) setupServer() {
	s.server = echo.New()

	s.server.Use(middleware.CORS())
	s.server.Use(middleware.Recover())
	s.server.Use(InitLogger(s.config))

	_ = s.CreateCloudGroup()

	s.server.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *ServerHttp) Start(_ context.Context) error {
	s.setupServer()
	return s.server.Start(s.config.Address)
}

func (s *ServerHttp) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
