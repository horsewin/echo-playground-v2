package model

type Reservation struct {
	PetId           string `json:"petId"`
	UserId          string `json:"user_id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	ReservationDate string `json:"reservation_date"`
}
