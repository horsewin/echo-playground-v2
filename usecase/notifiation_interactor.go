package usecase

import (
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
		app, err = interactor.NotificationRepository.Where(id)
		if err != nil {
			err = utils.SetErrorMassage("10001E")
			return
		}
	}

	return
}

// GetUnreadNotificationCount ...
func (interactor *NotificationInteractor) GetUnreadNotificationCount() (count model.NotificationCount, err error) {

	count, err = interactor.NotificationRepository.Count("unread = ?", true)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// MarkNotificationsRead ...
func (interactor *NotificationInteractor) MarkNotificationsRead() (notification model.Notification, err error) {
	notification, err = interactor.NotificationRepository.Update(map[string]interface{}{"Unread": false}, "unread = ?", true)

	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}
