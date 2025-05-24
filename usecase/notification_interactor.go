package usecase

import (
	"context"
	"fmt"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/utils"
)

// NotificationInteractor ...
type NotificationInteractor struct {
	NotificationRepository repository.NotificationRepositoryInterface
}

// GetNotifications ...
func (interactor *NotificationInteractor) GetNotifications(ctx context.Context, id string) (app model.Notifications, err error) {
	if id == "" {
		app, err = interactor.NotificationRepository.FindAll(ctx)
		if err != nil {
			err = utils.ConvertErrorMassage(ctx, "10001E", err)
			return
		}

	} else {
		app, err = interactor.NotificationRepository.Find(ctx, id)
		if err != nil {
			err = utils.ConvertErrorMassage(ctx, "10001E", err)
			return
		}
	}

	return
}

// GetUnreadNotificationCount ...
func (interactor *NotificationInteractor) GetUnreadNotificationCount(ctx context.Context) (count model.NotificationCount, err error) {
	whereClause := "is_read = :is_read"
	whereArgs := map[string]interface{}{"is_read": false}
	count, err = interactor.NotificationRepository.Count(ctx, whereClause, whereArgs)
	if err != nil {
		fmt.Println(err)
		err = utils.ConvertErrorMassage(ctx, "10001E", err)
		return
	}

	return
}

// MarkNotificationsRead ...
func (interactor *NotificationInteractor) MarkNotificationsRead(ctx context.Context) (err error) {
	clause := "unread = :unread"
	args := map[string]interface{}{"unread": true}
	err = interactor.NotificationRepository.Update(ctx, map[string]interface{}{"Unread": false}, clause, args)

	if err != nil {
		err = utils.ConvertErrorMassage(ctx, "10001E", err)
		return
	}

	return
}
