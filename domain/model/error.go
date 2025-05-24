package model

// ErrorMessages ... entity for error result
type ErrorMessages struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
