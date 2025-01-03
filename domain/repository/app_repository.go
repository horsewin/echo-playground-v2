package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

// AppRepositoryInterface ...
type AppRepositoryInterface interface {
	FindAll() (items model.Items, err error)
	Find(whereClause string, whereArgs map[string]interface{}) (items model.Items, err error)
	Create(input model.Item) (out model.Response, err error)
	Update(in map[string]interface{}, query string, args map[string]interface{}) (err error)
}

// AppRepository ...
type AppRepository struct {
	database.SQLHandler
}

// TODO: itemにする
const ItemsTable = "items"

// FindAll ...
func (repo *AppRepository) FindAll() (items model.Items, err error) {
	err = repo.SQLHandler.Scan(&items.Data, ItemsTable, "id desc")
	return items, err
}

// Find ...
func (repo *AppRepository) Find(whereClause string, whereArgs map[string]interface{}) (items model.Items, err error) {
	err = repo.SQLHandler.Where(&items.Data, ItemsTable, whereClause, whereArgs)
	return
}

// Create ...
func (repo *AppRepository) Create(input model.Item) (out model.Response, err error) {
	// inputをmap[string]interface{}に変換
	in := map[string]interface{}{
		"title":      input.Title,
		"name":       input.Name,
		"favorite":   input.Favorite,
		"img":        input.Img,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
	}

	err = repo.SQLHandler.Create(in, ItemsTable)

	if err != nil {
		return model.Response{
			Code:    400,
			Message: "Create error",
		}, err
	}

	return model.Response{
		Code:    200,
		Message: "OK",
	}, nil
}

// Update ...
func (repo *AppRepository) Update(in map[string]interface{}, query string, args map[string]interface{}) (err error) {
	err = repo.SQLHandler.Update(in, ItemsTable, query, args)
	return
}
