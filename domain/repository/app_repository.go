package repository

import (
	"fmt"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

// AppRepositoryInterface ...
type AppRepositoryInterface interface {
	FindAll() (items model.Items, err error)
	Find(query interface{}, args ...interface{}) (items model.Items, err error)
	Create(input model.Item) (out model.Response, err error)
	Update(value map[string]interface{}, query interface{}, args ...interface{}) (item model.Item, err error)
}

// AppRepository ...
type AppRepository struct {
	database.SQLHandler
}

// FindAll ...
func (repo *AppRepository) FindAll() (items model.Items, err error) {
	// TODO: impl
	return model.Items{}, fmt.Errorf("not implemented")
}

// Find ...
func (repo *AppRepository) Find(query interface{}, args ...interface{}) (items model.Items, err error) {
	//repo.SQLHandler.Where(&items.Data, query, args...)
	//return
	// TODO: impl
	return model.Items{}, fmt.Errorf("not implemented")
}

// Create ...
func (repo *AppRepository) Create(input model.Item) (out model.Response, err error) {
	// inputをmap[string]interface{}に変換
	in := map[string]interface{}{
		"id":         input.ID,
		"title":      input.Title,
		"name":       input.Name,
		"favorite":   input.Favorite,
		"img":        input.Img,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
	}

	err = repo.SQLHandler.Create(in)

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
func (repo *AppRepository) Update(value map[string]interface{}, query interface{}, args ...interface{}) (item model.Item, err error) {
	//repo.SQLHandler.Update(&item, value, query, args...)
	// TODO: impl
	return model.Item{}, fmt.Errorf("not implemented")
}
