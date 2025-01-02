package repository

import (
	"fmt"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

// NotificationRepositoryInterface ...
type NotificationRepositoryInterface interface {
	Where(id string) (account model.Notifications, err error)
	FindAll() (notifications model.Notifications, err error)
	Count(query interface{}, args ...interface{}) (data model.NotificationCount, err error)
	Update(value map[string]interface{}, query interface{}, args ...interface{}) (notification model.Notification, err error)
}

// NotificationRepository ....
type NotificationRepository struct {
	database.SQLHandler
}

// Where ...
func (repo *NotificationRepository) Where(id string) (notifications model.Notifications, err error) {
	repo.SQLHandler.Where(&notifications.Data, "id = ?", id)
	return
}

// FindAll ...
func (repo *NotificationRepository) FindAll() (notifications model.Notifications, err error) {
	// TODO: impl
	//repo.SQLHandler.Scan(&notifications.Data, "id desc")
	return model.Notifications{}, fmt.Errorf("not implemented")
}

// Count ...
func (repo *NotificationRepository) Count(query interface{}, args ...interface{}) (data model.NotificationCount, err error) {
	// TODO: impl
	return model.NotificationCount{Data: 0}, fmt.Errorf("not implemented")

	//var count int
	//repo.SQLHandler.Count(&count, &model.Notification{}, query, args...)
	//
	//return model.NotificationCount{Data: count}, nil
}

// Update ...
func (repo *NotificationRepository) Update(value map[string]interface{}, query interface{}, args ...interface{}) (notification model.Notification, err error) {
	// NOTE: When update with struct, GORM will only update non-zero fields,
	// you might want to use map to update attributes or use Select to specify fields to update
	//repo.SQLHandler.Update(&notification, value, query, args...)
	//return
	// TODO: impl
	return model.Notification{}, fmt.Errorf("not implemented")
}
