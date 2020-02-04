package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	echoWebServer *echo.Echo
	kube          *KubeMQ
	cfg           *Config
}

func NewServer(kube *KubeMQ, config *Config) (*Server, error) {
	s := &Server{
		echoWebServer: nil,
		kube:          kube,
		cfg:           config,
	}

	e := echo.New()
	s.echoWebServer = e

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.HTTPErrorHandler = s.customHTTPErrorHandler
	e.POST("/register", func(c echo.Context) error {
		return s.register(c)
	})
	e.POST("/verify", func(c echo.Context) error {
		return s.verify(c)
	})
	e.POST("/login", func(c echo.Context) error {
		return s.login(c)
	})
	e.POST("/logout", func(c echo.Context) error {
		return s.logout(c)
	})


	go func() {
		_ = s.echoWebServer.Start(":" + s.cfg.Port)
	}()

	return s, nil
}

type ErrorMessage struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (s *Server) customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		c.JSON(code, &ErrorMessage{
			ErrorCode: he.Code,
			Message:   he.Message.(string),
		})

	}
}
