package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

// PetRepositoryInterface ...
type PetRepositoryInterface interface {
	FindAll() (pets model.Pets, err error)
	Find(whereClause string, whereArgs map[string]interface{}) (pets model.Pets, err error)
	Create(input model.Pet) (out model.Response, err error)
	Update(in map[string]interface{}, query string, args map[string]interface{}) (err error)
}

// PetRepository ...
type PetRepository struct {
	database.SQLHandler
}

const PetsTable = "pets"

// FindAll ...
func (repo *PetRepository) FindAll() (pets model.Pets, err error) {
	err = repo.SQLHandler.Scan(&pets.Data, PetsTable, "id desc")
	return pets, err
}

// Find ...
func (repo *PetRepository) Find(whereClause string, whereArgs map[string]interface{}) (pets model.Pets, err error) {
	err = repo.SQLHandler.Where(&pets.Data, PetsTable, whereClause, whereArgs)
	return
}

// Create ...
func (repo *PetRepository) Create(input model.Pet) (out model.Response, err error) {
	// inputをmap[string]interface{}に変換
	in := map[string]interface{}{
		"id":               input.ID,
		"name":             input.Name,
		"breed":            input.Breed,
		"gender":           input.Gender,
		"price":            input.Price,
		"image_url":        input.ImageURL,
		"likes":            input.Likes,
		"shop_name":        input.ShopName,
		"shop_location":    input.ShopLocation,
		"birth_date":       input.BirthDate,
		"reference_number": input.ReferenceNumber,
		"tags":             input.Tags,
	}
	err = repo.SQLHandler.Create(in, PetsTable)

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
func (repo *PetRepository) Update(in map[string]interface{}, query string, args map[string]interface{}) (err error) {
	err = repo.SQLHandler.Update(in, PetsTable, query, args)
	return
}
