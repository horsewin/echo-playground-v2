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
