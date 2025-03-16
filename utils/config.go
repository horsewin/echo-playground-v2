package utils

import (
	"os"
	"strings"
)

// APIConfig ...
type APIConfig struct {
	HeaderValue struct {
		ClientID string
	}
	EnableTracing bool
}

// ConfigDB ...
type ConfigDB struct {
	Postgres struct {
		DBMS     string
		Username string
		Password string
		DBName   string
	}
}

// NewAPIConfig ...
func NewAPIConfig() *APIConfig {
	config := new(APIConfig)
	config.HeaderValue.ClientID = os.Getenv("SBCNTR_CLIENT_ID_HEADER")

	// 環境変数[SBCNTR_ENABLE_TRACING]を見てトレースを有効にする。対応しているTracingはAWS_XRAYのみ。
	// 環境変数[AWS_XRAY_SDK_DISABLED]がtrueの場合は必ずトレースを無効にする。
	enableKey := os.Getenv("SBCNTR_ENABLE_TRACING")
	if !SdkDisabled() && (strings.ToLower(enableKey) == "true" || enableKey == "1") {
		os.Setenv("AWS_XRAY_SDK_DISABLED", "FALSE")
		config.EnableTracing = true
	} else {
		os.Setenv("AWS_XRAY_SDK_DISABLED", "TRUE")
		config.EnableTracing = false
	}

	return config
}

// NewConfigDB ...
func NewConfigDB() *ConfigDB {
	config := new(ConfigDB)

	config.Postgres.DBMS = "postgres"
	config.Postgres.Username = os.Getenv("DB_USERNAME")
	config.Postgres.Password = os.Getenv("DB_PASSWORD")
	config.Postgres.DBName = os.Getenv("DB_NAME")

	return config
}

// Check if SDK is disabled
func SdkDisabled() bool {
	disableKey := os.Getenv("AWS_XRAY_SDK_DISABLED")
	return strings.ToLower(disableKey) == "true"
}
