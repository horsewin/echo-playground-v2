package model

// App ... entity for db result
type App struct {
	ID      string `json:"id" db:"id"`
	Message string `json:"message" db:"message"`
}

// Response ...
type Response struct {
	Code    int    `json:"code" xml:"code"`
	Message string `json:"msg" xml:"msg"`
}

// Hello ... entity for hello message
type Hello struct {
	Data string `json:"data" xml:"data"`
}

type (
	// Item ...
	Item struct {
		ID        int    `json:"id" db:"id"`
		Title     string `json:"title" db:"title"`
		Name      string `json:"name" db:"name"`
		Favorite  bool   `json:"favorite" db:"favorite"`
		Img       string `json:"img" db:"img"`
		CreatedAt string `json:"createdAt" db:"created_at"`
		UpdatedAt string `json:"updatedAt" db:"updated_at"`
	}

	// Items ...
	Items struct {
		Data []Item `json:"data"`
	}
)
