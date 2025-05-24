package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/horsewin/echo-playground-v2/domain/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
)

var messageConfig map[string]interface{}

func init() {
	err := json.Unmarshal([]byte(MessagesConfig), &messageConfig)
	if err != nil {
		panic(err)
	}
}

// GetMessageStatusCode is ...
func getMessageStatusCode(messageCode string) int {
	return int((messageConfig[messageCode].(map[string]interface{}))["statusCode"].(float64))
}

// getMessage is ...
func getMessage(lang string, messageCode string, args ...interface{}) string {
	llang := lang

	if !(llang == "ja" || llang == "en") {
		llang = "ja"
	}
	return fmt.Sprintf(((messageConfig[messageCode].(map[string]interface{}))["message"].(map[string]interface{}))[llang].(string), args...)
}

// GetError is ...
func GetError(context interface{}, lang string, err1 error) (err error) {
	c := context.(echo.Context)
	messageCode := err1.Error()

	// messageCodeが"数字5桁+E"の正規表現になっているかチェック
	if !regexp.MustCompile(`\d{5}[IWE]`).Match([]byte(messageCode)) {
		return c.JSON(http.StatusInternalServerError, &model.ErrorMessages{
			Code:    http.StatusInternalServerError,
			Message: "Unhandled internal server error",
		})
	}

	errStatusCode := getMessageStatusCode(messageCode)

	return echo.NewHTTPError(errStatusCode, getMessage(lang, messageCode))
}

// ConvertErrorMassage is ...
func ConvertErrorMassage(context interface{}, messageCode string, errorAsLogging error) (err error) {
	if errorAsLogging != nil && context != nil {
		// Try to use echo.Context logger if available
		if echoCtx, ok := context.(echo.Context); ok {
			echoCtx.Logger().Error(errorAsLogging)
		} else {
			// Fallback to standard output for other context types
			fmt.Printf("Error: %v\n", errorAsLogging)
		}
	}

	err = errors.New(messageCode)
	return
}
