#!/bin/bash

echo "🧪 APIテストを実行中..."
echo ""

# ベースURL
BASE_URL="http://localhost:8081"

echo "2️⃣ Hello World"
curl -s $BASE_URL/v1/helloworld | jq '.' || echo "Response: $(curl -s $BASE_URL/v1/helloworld)"
echo ""

echo "3️⃣ ペット一覧取得"
curl -s $BASE_URL/v1/pets | jq '.' || echo "Response: $(curl -s $BASE_URL/v1/pets)"
echo ""

echo "✅ テスト完了!"