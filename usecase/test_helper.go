package usecase

import (
	"context"
	"os"
)

// init はパッケージ初期化時にX-Rayを無効化する
func init() {
	// テスト実行時はX-Rayを無効化
	os.Setenv("AWS_XRAY_SDK_DISABLED", "true")
}

// testContext はテスト用のコンテキストを作成する
func testContext() context.Context {
	return context.Background()
}
