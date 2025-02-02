package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
	"time"
)

// ReservationRepositoryInterface ...
type ReservationRepositoryInterface interface {
	Create(input *model.Reservation) (err error)
}

// ReservationRepository ...
type ReservationRepository struct {
	database.SQLHandler
}

const ReservationTable = "reservations"

// Create ...
func (repo *ReservationRepository) Create(input *model.Reservation) (err error) {
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
	err = repo.SQLHandler.Create(in, ReservationTable)

	return
}
