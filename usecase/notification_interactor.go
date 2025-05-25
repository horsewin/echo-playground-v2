package usecase

import (
	"context"
	"fmt"

	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/model/errors"
	"github.com/horsewin/echo-playground-v2/domain/repository"
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
			return app, errors.NewBusinessError("10001E", err)
		}

	} else {
		app, err = interactor.NotificationRepository.Find(ctx, id)
		if err != nil {
			return app, errors.NewBusinessError("10001E", err)
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
		return count, errors.NewBusinessError("10001E", err)
	}

	return
}

// MarkNotificationsRead ...
func (interactor *NotificationInteractor) MarkNotificationsRead(ctx context.Context) (err error) {
	clause := "is_read = :is_read"
	args := map[string]interface{}{"is_read": false}
	err = interactor.NotificationRepository.Update(ctx, map[string]interface{}{"is_read": true}, clause, args)

	if err != nil {
		return errors.NewBusinessError("10001E", err)
	}

	return
}
