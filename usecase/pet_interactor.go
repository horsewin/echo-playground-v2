package usecase

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/domain/repository"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/jinzhu/copier"
)

// PetInteractor ...
type PetInteractor struct {
	PetRepository         repository.PetRepositoryInterface
	ReservationRepository repository.ReservationRepositoryInterface
	FavoriteRepository    repository.FavoriteRepositoryInterface
}

// GetPets ...
func (interactor *PetInteractor) GetPets(filter *model.PetFilter) (pets []model.Pet, err error) {
	// Repository層からデータを取得
	_app, err := interactor.PetRepository.Find(filter)
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
	// like状態を取得
	favMap, err := interactor.FavoriteRepository.FindByUserId(input.UserId)
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}
	if favMap[input.PetId].Value == input.Value {
		err = utils.SetErrorMassage("00001I")
		return
	}

	// 現在のペットモデルを取得
	petData, err := interactor.PetRepository.Find(&model.PetFilter{ID: input.PetId})
	if err != nil {
		err = utils.SetErrorMassage("10001E")
		return
	}
	var pet model.Pet
	for _, p := range petData.Data {
		if p.ID == input.PetId {
			err = copier.Copy(&pet, &p)
			if err != nil {
				err = utils.SetErrorMassage("10002E")
				return
			}

			// マッピングしきれない構造は手動でコピー
			pet.Shop = model.Shop{
				Name:     p.ShopName,
				Location: p.ShopLocation,
			}
		}
	}

	// Like数を更新
	if input.Value {
		pet.Likes = pet.Likes + 1
	} else {
		pet.Likes = pet.Likes - 1
	}

	// TODO: トランザクションを使って更新処理を行う
	err = interactor.PetRepository.Update(&pet)
	if err != nil {
		err = utils.SetErrorMassage("10003E")
		return
	}

	if input.Value {
		err = interactor.FavoriteRepository.Create(&model.Favorite{
			PetId:  input.PetId,
			UserId: input.UserId,
			Value:  input.Value,
		})
		if err != nil {
			err = utils.SetErrorMassage("10004E")
			return
		}
	} else {
		err = interactor.FavoriteRepository.Delete(&model.Favorite{
			PetId:  input.PetId,
			UserId: input.UserId,
		})
		if err != nil {
			err = utils.SetErrorMassage("10005E")
			return
		}
	}

	return
}

// CreateReservation ...
func (interactor *PetInteractor) CreateReservation(input *model.Reservation) (err error) {
	err = interactor.ReservationRepository.Create(input)
	if err != nil {
		err = utils.SetErrorMassage("10003E")
		return
	}
	return
}
