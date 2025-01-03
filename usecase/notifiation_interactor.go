package usecase

import (
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
func (interactor *NotificationInteractor) GetNotifications(id string) (app model.Notifications, err error) {
	if id == "" {
		app, err = interactor.NotificationRepository.FindAll()
		if err != nil {
			err = utils.SetErrorMassage("10001E")
			return
		}

	} else {
		app, err = interactor.NotificationRepository.Find(id)
		if err != nil {
			err = utils.SetErrorMassage("10001E")
			return
		}
	}

	return
}

// GetUnreadNotificationCount ...
func (interactor *NotificationInteractor) GetUnreadNotificationCount() (count model.NotificationCount, err error) {
	whereClause := "unread = :unread"
	whereArgs := map[string]interface{}{"unread": true}
	count, err = interactor.NotificationRepository.Count(whereClause, whereArgs)
	if err != nil {
		fmt.Println(err)
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// MarkNotificationsRead ...
func (interactor *NotificationInteractor) MarkNotificationsRead() (err error) {
	clause := "unread = :unread"
	args := map[string]interface{}{"unread": true}
	err = interactor.NotificationRepository.Update(map[string]interface{}{"Unread": false}, clause, args)

	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}
