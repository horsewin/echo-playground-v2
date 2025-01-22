package usecase

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/utils"
)

// PetInteractor ...
type PetInteractor struct {
	PetRepository repository.PetRepositoryInterface
}

// Getpets ...
func (interactor *PetInteractor) Getpets(gender string) (app model.Pets, err error) {
	var query string
	var args map[string]interface{}
	if gender == "male" {
		query = "gender = :gender"
		args = map[string]interface{}{"gender": "male"}
	} else if gender == "female" {
		query = "gender = :gender"
		args = map[string]interface{}{"gender": "female"}
	} else {
		query = ""
		args = map[string]interface{}{}
	}

	app, err = interactor.PetRepository.Find(query, args)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// CreateItem ...
func (interactor *PetInteractor) CreateItem(input model.Pet) (response model.Response, err error) {
	response, err = interactor.PetRepository.Create(input)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

// UpdateFavoriteAttr ...
func (interactor *PetInteractor) UpdateFavoriteAttr(input model.Pet) (err error) {
	whereClause := "id = :id"
	whereArgs := map[string]interface{}{"id": input.ID}
	err = interactor.PetRepository.Update(map[string]interface{}{"Favorite": input.Likes}, whereClause, whereArgs)

	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}
