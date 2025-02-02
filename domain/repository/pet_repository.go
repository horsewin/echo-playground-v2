package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/lib/pq"
	"time"
)

// pet... pets テーブルの各カラムと対応する構造体
type pet struct {
	ID              string         `db:"id"`
	Name            string         `db:"name"`
	Breed           string         `db:"breed"`
	Gender          string         `db:"gender"`
	Price           float64        `db:"price"`
	ImageURL        *string        `db:"image_url"`
	Likes           int            `db:"likes"`
	ShopName        string         `db:"shop_name"`
	ShopLocation    string         `db:"shop_location"`
	BirthDate       *time.Time     `db:"birth_date"`
	ReferenceNumber string         `db:"reference_number"`
	Tags            pq.StringArray `db:"tags"`
	CreatedAt       *time.Time     `db:"created_at"`
	UpdatedAt       *time.Time     `db:"updated_at"`
}

type pets struct {
	Data []pet
}

// PetRepositoryInterface ...
type PetRepositoryInterface interface {
	FindAll() (pets pets, err error)
	Find(whereClause string, whereArgs map[string]interface{}) (pets pets, err error)
	Update(input *model.Pet) (err error)
}

// PetRepository ...
type PetRepository struct {
	database.SQLHandler
}

const PetsTable = "pets"

// FindAll ...
func (repo *PetRepository) FindAll() (pets pets, err error) {
	err = repo.SQLHandler.Scan(&pets.Data, PetsTable, "id desc")
	return pets, err
}

// Find ...
func (repo *PetRepository) Find(whereClause string, args map[string]interface{}) (pets pets, err error) {
	err = repo.SQLHandler.Where(&pets.Data, PetsTable, whereClause, args)
	return
}

// Update ...
func (repo *PetRepository) Update(input *model.Pet) (err error) {
	// Petドメインモデルをリポジトリモデルに変換
	now := time.Now()
	in := map[string]interface{}{
		"id":               input.ID,
		"name":             input.Name,
		"breed":            input.Breed,
		"gender":           input.Gender,
		"price":            input.Price,
		"image_url":        input.ImageURL,
		"likes":            input.Likes,
		"shop_name":        input.Shop.Name,
		"shop_location":    input.Shop.Location,
		"birth_date":       input.BirthDate,
		"reference_number": input.ReferenceNumber,
		"tags":             pq.StringArray(input.Tags),
		"updated_at":       &now,
	}

	// クエリ組み立て
	whereClause := "id = :id"

	err = repo.SQLHandler.Update(in, PetsTable, whereClause)
	return
}
