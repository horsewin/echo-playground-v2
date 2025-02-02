package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

type favorite struct {
	ID     string `db:"id"`
	PetId  string `db:"pet_id"`
	UserId string `db:"user_id"`
}

type favorites struct {
	Data []favorite
}

// FavoriteRepositoryInterface ...
type FavoriteRepositoryInterface interface {
	FindByUserId(userId string) (favorites map[string]model.Favorite, err error)
	Create(input *model.Favorite) (err error)
	Delete(input *model.Favorite) (err error)
}

// FavoriteRepository ...
type FavoriteRepository struct {
	database.SQLHandler
}

const FavoriteTable = "favorites"

// FindByUserId ...
func (f FavoriteRepository) FindByUserId(userId string) (favMap map[string]model.Favorite, err error) {
	// inputをmapに変換
	in := map[string]interface{}{"user_id": userId}

	// リポジトリモデルをDBから取得
	var _favorites favorites
	err = f.SQLHandler.Where(&_favorites.Data, FavoriteTable, "user_id = :user_id", in)
	if err != nil {
		return
	}

	// ドメインモデルに変換
	favMap = make(map[string]model.Favorite)
	for _, _favorite := range _favorites.Data {
		favMap[_favorite.PetId] = model.Favorite{
			Id:     _favorite.ID,
			PetId:  _favorite.PetId,
			UserId: _favorite.UserId,
			Value:  true,
		}
	}

	return
}

// Create ...
func (f FavoriteRepository) Create(input *model.Favorite) (err error) {
	// ドメインモデルをmapに変換
	in := map[string]interface{}{"pet_id": input.PetId, "user_id": input.UserId}

	// リポジトリモデルをDBに保存
	err = f.SQLHandler.Create(in, FavoriteTable)

	return
}

// Delete ...
func (f FavoriteRepository) Delete(input *model.Favorite) (err error) {
	// ドメインモデルをリポジトリモデルに変換
	_favorite := favorite{
		PetId:  input.PetId,
		UserId: input.UserId,
	}

	// リポジトリモデルをmapに変換
	in := map[string]interface{}{"pet_id": _favorite.PetId, "user_id": _favorite.UserId}

	// リポジトリモデルをDBから削除
	err = f.SQLHandler.Delete(in, FavoriteTable)

	return
}
