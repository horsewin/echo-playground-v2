package handlers

import (
	"net/http"

	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/labstack/echo/v4"
)

// HelloWorldHandler ...
type HelloWorldHandler struct {
}

// NewHelloWorldHandler ...
func NewHelloWorldHandler() *HelloWorldHandler {
	return &HelloWorldHandler{}
}

// SayHelloWorld ...
func (handler *HelloWorldHandler) SayHelloWorld() echo.HandlerFunc {
	body := &model.Hello{
		Data: "Hello world",
	}
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, body)
	}
}
