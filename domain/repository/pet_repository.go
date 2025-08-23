package repository

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/lib/pq"
)

const PetsTable = "pets"

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
	Find(ctx context.Context, filter *model.PetFilter) (pets pets, err error)
	Update(ctx context.Context, input *model.Pet) (err error)
}

// PetRepository ...
type PetRepository struct {
	database.SQLHandler
}

// Find ...
func (repo *PetRepository) Find(ctx context.Context, filter *model.PetFilter) (pets pets, err error) {
	// サブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "PetRepository.Find")
	defer func() {
		if seg != nil {
			seg.Close(err)
		}
	}()

	// フィルタ条件をリポジトリで解釈する型に変換
	whereClause, args := parseFilter(filter)

	// メタデータを追加
	if err := seg.AddMetadata("filter", filter); err != nil {
		utils.LogError("Failed to add filter metadata: %v", err)
	}
	if err := seg.AddMetadata("where_clause", whereClause); err != nil {
		utils.LogError("Failed to add where_clause metadata: %v", err)
	}
	if err := seg.AddMetadata("args", args); err != nil {
		utils.LogError("Failed to add args metadata: %v", err)
	}

	// インフラストラクチャレイヤの処理を実行
	err = repo.SQLHandler.Where(subCtx, &pets.Data, PetsTable, strings.Join(whereClause, " and "), args)

	// 結果のメタデータを追加
	if err == nil {
		if err := seg.AddMetadata("result_count", len(pets.Data)); err != nil {
			utils.LogError("Failed to add result_count metadata: %v", err)
		}
	}

	return
}

// Update ...
func (repo *PetRepository) Update(ctx context.Context, input *model.Pet) (err error) {
	// サブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "PetRepository.Update")
	defer func() {
		if seg != nil {
			seg.Close(err)
		}
	}()

	// メタデータを追加
	if err := seg.AddMetadata("pet_id", input.ID); err != nil {
		utils.LogError("Failed to add pet_id metadata: %v", err)
	}

	// Petドメインモデルをリポジトリモデルに変換
	now := time.Now()
	setParams := map[string]interface{}{
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

	// WHERE句のパラメータ
	whereClause := "id = :id"
	whereParams := map[string]interface{}{
		"id": input.ID,
	}

	// SQLHandlerを呼び出し
	err = repo.SQLHandler.Update(subCtx, setParams, PetsTable, whereClause, whereParams)

	return
}

// parseFilter ... フィルタ条件を解釈してクエリ条件とバインド変数を返す
func parseFilter(filter *model.PetFilter) ([]string, map[string]interface{}) {
	args := map[string]interface{}{}
	whereClause := make([]string, 0)

	if filter != nil {
		if filter.Gender != "" {
			if strings.EqualFold(filter.Gender, "male") {
				whereClause = append(whereClause, "gender = :gender")
				args["gender"] = "Male"
			} else if strings.EqualFold(filter.Gender, "female") {
				whereClause = append(whereClause, "gender = :gender")
				args["gender"] = "Female"
			}
		}
		if filter.Price != 0 {
			whereClause = append(whereClause, "price = :price")
			args["price"] = filter.Price
		}
		if filter.Name != "" {
			whereClause = append(whereClause, "name = :name")
			args["name"] = filter.Name
		}
		if filter.ID != "" {
			whereClause = append(whereClause, "id = :id")
			args["id"] = filter.ID
		}
		if filter.ReferenceNumber != "" {
			whereClause = append(whereClause, "reference_number = :reference_number")
			args["reference_number"] = filter.ReferenceNumber
		}
		if filter.Breed != "" {
			whereClause = append(whereClause, "breed = :breed")
			args["breed"] = filter.Breed
		}
	}
	return whereClause, args
}
