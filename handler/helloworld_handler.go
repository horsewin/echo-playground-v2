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
		Message: "Hello world",
	}

	// ドメインモデルをJSONにして返却
	json := model.APIResponse{
		Data: body,
	}

	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, json)
	}
}
