package utils

import (
	"fmt"
	"log"
)

// LogError はエラーメッセージをログに記録します
func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// LogInfo は情報メッセージをログに記録します
func LogInfo(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

// LogWarning は警告メッセージをログに記録します
func LogWarning(format string, args ...interface{}) {
	log.Printf("[WARNING] "+format, args...)
}

// LogDebug はデバッグメッセージをログに記録します
func LogDebug(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

// FormatError はエラーメッセージをフォーマットして返します
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}
