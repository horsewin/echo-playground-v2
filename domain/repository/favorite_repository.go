package repository

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/utils"
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
	FindByUserId(ctx context.Context, userId string) (favorites map[string]model.Favorite, err error)
	Create(ctx context.Context, input *model.Favorite) (err error)
	Delete(ctx context.Context, input *model.Favorite) (err error)
}

// FavoriteRepository ...
type FavoriteRepository struct {
	database.SQLHandler
}

const FavoriteTable = "favorites"

// FindByUserId ...
func (f FavoriteRepository) FindByUserId(ctx context.Context, userId string) (favMap map[string]model.Favorite, err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "FavoriteRepository.FindByUserId")
	defer func() {
		if seg != nil {
			seg.Close(err)
		}
	}()

	// メタデータを追加
	if err := seg.AddMetadata("user_id", userId); err != nil {
		utils.LogError("Failed to add user_id metadata: %v", err)
	}

	// inputをmapに変換
	in := map[string]interface{}{"user_id": userId}

	// リポジトリモデルをDBから取得
	var _favorites favorites
	err = f.SQLHandler.Where(ctx, &_favorites.Data, FavoriteTable, "user_id = :user_id", in)
	if err != nil {
		return
	}

	// 結果のメタデータを追加
	if err := seg.AddMetadata("result_count", len(_favorites.Data)); err != nil {
		utils.LogError("Failed to add result_count metadata: %v", err)
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
func (f FavoriteRepository) Create(ctx context.Context, input *model.Favorite) (err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "FavoriteRepository.Create")
	defer func() {
		if seg != nil {
			seg.Close(err)
		}
	}()

	// メタデータを追加
	if err := seg.AddMetadata("pet_id", input.PetId); err != nil {
		utils.LogError("Failed to add pet_id metadata: %v", err)
	}
	if err := seg.AddMetadata("user_id", input.UserId); err != nil {
		utils.LogError("Failed to add user_id metadata: %v", err)
	}
	if err := seg.AddMetadata("value", input.Value); err != nil {
		utils.LogError("Failed to add value metadata: %v", err)
	}

	// ドメインモデルをmapに変換
	in := map[string]interface{}{"pet_id": input.PetId, "user_id": input.UserId}

	// リポジトリモデルをDBに保存
	err = f.SQLHandler.Create(ctx, in, FavoriteTable)

	return
}

// Delete ...
func (f FavoriteRepository) Delete(ctx context.Context, input *model.Favorite) (err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "FavoriteRepository.Delete")
	defer func() {
		if seg != nil {
			seg.Close(err)
		}
	}()

	// メタデータを追加
	if err := seg.AddMetadata("pet_id", input.PetId); err != nil {
		utils.LogError("Failed to add pet_id metadata: %v", err)
	}
	if err := seg.AddMetadata("user_id", input.UserId); err != nil {
		utils.LogError("Failed to add user_id metadata: %v", err)
	}

	// ドメインモデルをリポジトリモデルに変換
	_favorite := favorite{
		PetId:  input.PetId,
		UserId: input.UserId,
	}

	// リポジトリモデルをmapに変換
	in := map[string]interface{}{"pet_id": _favorite.PetId, "user_id": _favorite.UserId}

	// リポジトリモデルをDBから削除
	err = f.SQLHandler.Delete(ctx, in, FavoriteTable)

	return
}
