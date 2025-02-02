package repository

import (
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/horsewin/echo-playground-v2/interface/database"
)

type reservation struct {
	ID              string `db:"id"`
	PetId           string `db:"pet_id"`
	UserId          string `db:"user_id"`
	Email           string `db:"email"`
	FullName        string `db:"full_name"`
	ReservationDate string `db:"reservation_date"`
}

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
	// ドメインモデルをリポジトリモデルに変換
	_reservation := reservation{
		PetId:           input.PetId,
		UserId:          input.UserId,
		Email:           input.Email,
		FullName:        input.FullName,
		ReservationDate: input.ReservationDate,
	}

	// リポジトリモデルをmapに変換
	in := map[string]interface{}{"pet_id": _reservation.PetId, "user_id": _reservation.UserId, "email": _reservation.Email, "full_name": _reservation.FullName, "reservation_date": _reservation.ReservationDate}

	// リポジトリモデルをDBに保存
	err = repo.SQLHandler.Create(in, ReservationTable)

	return
}
