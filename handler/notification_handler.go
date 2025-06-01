package handlers

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/model/errors"

	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/usecase"
	"github.com/labstack/echo/v4"
)

// NotificationReadRequest JSON形式のリクエストボディをバインドするための構造体
type NotificationReadRequest struct {
	ID string `json:"id" form:"id"`
}

// NotificationHandler ...
type NotificationHandler struct {
	Interactor usecase.NotificationInteractor
}

// NewNotificationHandler ...
func NewNotificationHandler(sqlHandler database.SQLHandler) *NotificationHandler {
	return &NotificationHandler{
		Interactor: usecase.NotificationInteractor{
			NotificationRepository: &repository.NotificationRepository{
				SQLHandler: sqlHandler,
			},
		},
	}
}

// GetNotifications ...
func (handler *NotificationHandler) GetNotifications() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// サブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "GetNotifications")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			c.Logger().Warn("Failed to begin subsegment: GetNotifications")

			// contextを渡す（セグメントなし）
			resJSON, err := handler.Interactor.GetNotifications(ctx, c.QueryParam("id"))
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, resJSON)
		}
		defer seg.Close(err)

		id := c.QueryParam("id")

		// Add metadata to the segment
		if err := seg.AddMetadata("id", id); err != nil {
			c.Logger().Errorf("Failed to add id metadata: %v", err)
		}

		// contextを渡す
		resJSON, err := handler.Interactor.GetNotifications(ctx, id)
		if err != nil {
			return errors.NewEchoHTTPError(ctx, err)
		}

		return c.JSON(http.StatusOK, resJSON)
	}
}

// PostNotificationsRead ...
func (handler *NotificationHandler) PostNotificationsRead() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// サブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "PostNotificationsRead")
		if seg == nil {
			// セグメントが作成できない場合はログに記録して処理を続行
			logger := zerolog.Ctx(ctx)
			logger.Warn().Msg("Failed to begin subsegment: PostNotificationsRead")

			// JSONリクエストボディからnotificationIdを取得
			var req NotificationReadRequest
			if err := c.Bind(&req); err != nil {
				return errors.NewEchoHTTPError(ctx, err)
			}
			notificationId := req.ID
			fmt.Println("notificationId: ", notificationId)

			// contextを渡す（セグメントなし）
			err = handler.Interactor.MarkNotificationsRead(ctx, notificationId)
			if err != nil {
				return errors.NewEchoHTTPError(ctx, err)
			}

			return c.JSON(http.StatusOK, model.Response{
				Code:    http.StatusOK,
				Message: "OK",
			})
		}
		defer seg.Close(err)

		// JSONリクエストボディからnotificationIdを取得
		var req NotificationReadRequest
		if err := c.Bind(&req); err != nil {
			return errors.NewEchoHTTPError(ctx, err)
		}
		notificationId := req.ID
		fmt.Println("notificationId: ", notificationId)

		// Add metadata to the segment
		if err := seg.AddMetadata("id", notificationId); err != nil {
			c.Logger().Errorf("Failed to add id metadata: %v", err)
		}

		// contextを渡す
		err = handler.Interactor.MarkNotificationsRead(ctx, notificationId)
		if err != nil {
			return errors.NewEchoHTTPError(ctx, err)
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "OK",
		})
	}
}
