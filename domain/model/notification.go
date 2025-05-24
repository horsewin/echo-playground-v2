package model

type (
	// Notification ... entity for notification db result
	Notification struct {
		ID        int    `json:"id" db:"id"`
		Title     string `json:"title" db:"title"`
		UserId    string `json:"user_id" db:"user_id"`
		Message   string `json:"message" db:"message"`
		IsRead    bool   `json:"is_read" db:"is_read"`
		CreatedAt string `json:"createdAt" db:"created_at"`
		UpdatedAt string `json:"updatedAt" db:"updated_at"`
	}
	// Notifications ... array entity for notification
	Notifications struct {
		Data []Notification `json:"data"`
	}

	// NotificationCount ... unread count of notifications
	NotificationCount struct {
		Data int `json:"data"`
	}
)

// TableName ... override table name accessor
func (Notification) TableName() string {
	return "notifications"
}
