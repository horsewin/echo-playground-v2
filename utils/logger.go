package utils

import (
	"log"
)

// LogError はエラーメッセージをログに記録します
func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}
