#!/bin/bash

echo "🚀 アプリケーションを起動中..."

# 環境変数を設定
export DB_HOST=localhost
export DB_PORT=5433  # Docker Composeでマップされたポート
export DB_USERNAME=sbcntrapp
export DB_PASSWORD=password
export DB_NAME=sbcntrapp
export DB_CONN=1
export SBCNTR_ENABLE_TRACING=true
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

echo "📝 環境変数:"
echo "  DB_HOST=$DB_HOST:$DB_PORT"
echo "  DB_NAME=$DB_NAME"
echo "  OTEL_EXPORTER_OTLP_ENDPOINT=$OTEL_EXPORTER_OTLP_ENDPOINT"
echo ""

# アプリケーションを起動
echo "🏃 アプリケーション起動..."
go run main.go