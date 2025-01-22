package model

import (
	"github.com/lib/pq"
	"time"
)

// Hello ... entity for hello message
type Hello struct {
	Message string `json:"message" xml:"message"`
}

// Pet... pets テーブルの各カラムと対応する構造体
type Pet struct {
	ID              string         `json:"id" db:"id"`
	Name            string         `json:"name"  db:"name"`
	Breed           string         `json:"breed" db:"breed"`
	Gender          string         `json:"gender" db:"gender"`
	Price           float64        `json:"price" db:"price"`
	ImageURL        *string        `json:"image_url" db:"image_url"`
	Likes           int            `json:"likes" db:"likes"`
	ShopName        string         `json:"shop_name"      db:"shop_name"`
	ShopLocation    string         `json:"shop_location"  db:"shop_location"`
	BirthDate       *time.Time     `json:"birth_date" db:"birth_date"`
	ReferenceNumber string         `json:"reference_number" db:"reference_number"`
	Tags            pq.StringArray `json:"tags" db:"tags"`
	CreatedAt       *time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at" db:"updated_at"`
}

// 複数のペットをまとめる例
type Pets struct {
	Data []Pet `json:"data"`
}
