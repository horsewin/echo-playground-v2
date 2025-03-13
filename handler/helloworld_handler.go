package handlers

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"

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
		// ルーターで作成されたセグメントを使用するためにBeginSegmentは不要
		// 代わりにサブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "SayHelloWorld")
		defer seg.Close(nil) // エラーがない場合はnilを渡す

		// Add metadata to the segment
		if err := seg.AddMetadata("message", body.Message); err != nil {
			c.Logger().Errorf("Failed to add message metadata: %v", err)
		}

		return c.JSON(http.StatusOK, model.APIResponse{
			Data: body,
		})
	}
}
