package usecase

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/utils"
)

// AppInteractor ...
type AppInteractor struct {
	AppRepository repository.AppRepositoryInterface
}

// GetItems ...
func (interactor *AppInteractor) GetItems(favorite string) (app model.Items, err error) {
	var query string
	var args interface{}
	if favorite == "true" {
		query = "favorite = ?"
		args = true
	} else if favorite == "false" {
		query = "favorite = ?"
		args = false
	} else {
		query = ""
		args = ""
	}

	app, err = interactor.AppRepository.Find(query, args)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// CreateItem ...
func (interactor *AppInteractor) CreateItem(input model.Item) (response model.Response, err error) {
	response, err = interactor.AppRepository.Create(input)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// UpdateFavoriteAttr ...
func (interactor *AppInteractor) UpdateFavoriteAttr(input model.Item) (item model.Item, err error) {
	whereClause := "id = ?"
	whereArgs := map[string]interface{}{"id": input.ID}
	err = interactor.AppRepository.Update(map[string]interface{}{"Favorite": input.Favorite}, whereClause, whereArgs)

	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}
