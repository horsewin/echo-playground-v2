package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

// NotificationRepositoryInterface ...
type NotificationRepositoryInterface interface {
	Where(id string) (account model.Notifications, err error)
	FindAll() (notifications model.Notifications, err error)
	Count(query string, args map[string]interface{}) (data model.NotificationCount, err error)
	Update(in map[string]interface{}, query string, args map[string]interface{}) (err error)
}

// NotificationRepository ....
type NotificationRepository struct {
	database.SQLHandler
}

const TABLE_NAME = "Notification"

// Where ...
func (repo *NotificationRepository) Where(id string) (notifications model.Notifications, err error) {
	whereClause := "id = :id"
	whereArgs := map[string]interface{}{"id": id}
	err = repo.SQLHandler.Where(&notifications.Data, TABLE_NAME, whereClause, whereArgs)
	return
}

// FindAll ...
func (repo *NotificationRepository) FindAll() (notifications model.Notifications, err error) {
	err = repo.SQLHandler.Scan(&notifications.Data, TABLE_NAME, "id desc")
	return notifications, err
}

// Count ...
func (repo *NotificationRepository) Count(query string, args map[string]interface{}) (data model.NotificationCount, err error) {
	var count int
	err = repo.SQLHandler.Count(&count, TABLE_NAME, query, args)
	return model.NotificationCount{Data: count}, err
}

// Update ...
func (repo *NotificationRepository) Update(in map[string]interface{}, query string, args map[string]interface{}) (err error) {
	err = repo.SQLHandler.Update(in, TABLE_NAME, query, args)
	return
}
