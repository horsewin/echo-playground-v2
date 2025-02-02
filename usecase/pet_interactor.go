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

// GetPets ...
func (interactor *PetInteractor) GetPets(gender string) (pets []model.Pet, err error) {
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

	// Repository層からデータを取得
	_app, err := interactor.PetRepository.Find(query, args)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	// ドメインモデルに変換
	for _, p := range _app.Data {
		pets = append(pets, model.Pet{
			ID:       p.ID,
			Name:     p.Name,
			Breed:    p.Breed,
			Gender:   p.Gender,
			Price:    p.Price,
			ImageURL: p.ImageURL,
			Likes:    p.Likes,
			Shop: model.Shop{
				Name:     p.ShopName,
				Location: p.ShopLocation,
			},
			BirthDate:       p.BirthDate,
			ReferenceNumber: p.ReferenceNumber,
			Tags:            p.Tags,
		})
	}

	return pets, nil
}

// UpdateLikeCount ...
func (interactor *PetInteractor) UpdateLikeCount(input *model.InputUpdateLikeRequest) (err error) {
	whereClause := "id = :id"
	whereArgs := map[string]interface{}{"id": input.PetId}
	err = interactor.PetRepository.Update(map[string]interface{}{"Favorite": input.Likes}, whereClause, whereArgs)

	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}

	return
}

func (interactor *PetInteractor) CreateReservation(input *model.Reservation) (err error) {

}
