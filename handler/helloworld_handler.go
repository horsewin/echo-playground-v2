package handlers

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rs/zerolog"

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
	return func(c echo.Context) error {
		// Get logger from context
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)

		// サブセグメントを作成
		_, seg := xray.BeginSubsegment(ctx, "SayHelloWorld")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			logger.Warn().Msg("Failed to begin subsegment: SayHelloWorld")
			return c.JSON(http.StatusOK, model.APIResponse{
				Data: body,
			})
		}
		defer seg.Close(nil) // エラーがない場合はnilを渡す

		// Add metadata to the segment
		if err := seg.AddMetadata("message", body.Message); err != nil {
			logger.Error().Err(err).Msg("Failed to add message metadata")
		}

		return c.JSON(http.StatusOK, model.APIResponse{
			Data: body,
		})
	}
}

// SayError ...
func (handler *HelloWorldHandler) SayError() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Return a client-friendly error using echo.NewHTTPError
		// This will be handled by our custom error handler
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
}
