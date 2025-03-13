package handlers

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"

	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/usecase"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/labstack/echo/v4"
)

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
		// ルーターで作成されたセグメントを使用するためにBeginSegmentは不要
		// 代わりにサブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "GetNotifications")
		defer seg.Close(err)

		id := c.QueryParam("id")

		// Add metadata to the segment
		if err := seg.AddMetadata("id", id); err != nil {
			c.Logger().Errorf("Failed to add id metadata: %v", err)
		}

		// contextを渡す
		resJSON, err := handler.Interactor.GetNotifications(ctx, id)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, resJSON)
	}
}

// GetUnreadNotificationCount ...
func (handler *NotificationHandler) GetUnreadNotificationCount() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// ルーターで作成されたセグメントを使用するためにBeginSegmentは不要
		// 代わりにサブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "GetUnreadNotificationCount")
		defer seg.Close(err)

		// contextを渡す
		resJSON, err := handler.Interactor.GetUnreadNotificationCount(ctx)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, resJSON)
	}
}

// PostNotificationsRead ...
func (handler *NotificationHandler) PostNotificationsRead() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// ルーターで作成されたセグメントを使用するためにBeginSegmentは不要
		// 代わりにサブセグメントを作成
		ctx := c.Request().Context()
		_, seg := xray.BeginSubsegment(ctx, "PostNotificationsRead")
		defer seg.Close(err)

		// contextを渡す
		err = handler.Interactor.MarkNotificationsRead(ctx)
		if err != nil {
			return utils.GetErrorMassage(c, "en", err)
		}

		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "OK",
		})
	}
}
