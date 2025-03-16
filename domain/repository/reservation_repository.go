package repository

import (
	"context"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"github.com/horsewin/echo-playground-v2/utils"
)

// ReservationRepositoryInterface ...
type ReservationRepositoryInterface interface {
	Create(ctx context.Context, input *model.Reservation) (err error)
	GetCountByPetID(ctx context.Context, petID string) (count int64, err error)
}

// ReservationRepository ...
type ReservationRepository struct {
	database.SQLHandler
}

const ReservationTable = "reservations"

// Create ...
func (repo *ReservationRepository) Create(ctx context.Context, input *model.Reservation) (err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "ReservationRepository.Create")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("pet_id", input.PetId); err != nil {
		utils.LogError("Failed to add pet_id metadata: %v", err)
	}
	if err := seg.AddMetadata("user_id", input.UserId); err != nil {
		utils.LogError("Failed to add user_id metadata: %v", err)
	}
	if err := seg.AddMetadata("reservation_date", input.ReservationDate); err != nil {
		utils.LogError("Failed to add reservation_date metadata: %v", err)
	}

	// ドメインモデルをmapに変換
	rsvDatetime, err := time.Parse("20060102", input.ReservationDate)
	if err != nil {
		return
	}

	in := map[string]interface{}{
		"pet_id":    input.PetId,
		"user_id":   input.UserId,
		"email":     input.Email,
		"user_name": input.FullName,
		// yyyymmdd形式をdatetimeに変換
		"reservation_datetime": rsvDatetime,
	}

	// リポジトリモデルをDBに保存
	err = repo.SQLHandler.Create(ctx, in, ReservationTable)

	return
}

// GetCountByPetID ...
func (repo *ReservationRepository) GetCountByPetID(ctx context.Context, petID string) (count int64, err error) {
	// サブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "ReservationRepository.GetCountByPetID")
	defer seg.Close(err)

	// メタデータを追加
	if err := seg.AddMetadata("pet_id", petID); err != nil {
		utils.LogError("Failed to add pet_id metadata: %v", err)
	}

	var tmpCount int
	whereClause := "pet_id = :pet_id"
	whereArgs := map[string]interface{}{"pet_id": petID}
	err = repo.SQLHandler.Count(ctx, &tmpCount, ReservationTable, whereClause, whereArgs)
	count = int64(tmpCount)

	return
}
