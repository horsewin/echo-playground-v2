# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Echo v4フレームワークを使用したGo言語のREST APIサーバーです。クリーンアーキテクチャの原則に従ったペットショップのバックエンドAPIのデモアプリケーションです。

## 開発コマンド

### ビルドと実行
```bash
# 完全なワークフロー: validate → build → run
make all

# ビルドのみ
make build

# コンパイル済みバイナリの実行
make run
```

### コード品質チェック
```bash
# すべての検証チェックを実行 (fmt, vet, golangci-lint)
make validate

# テストの実行（注意：現在テストファイルは存在しません）
make test
```

### 依存関係管理
```bash
# 依存関係の更新
make update-deps

# Linux ARM64向けクロスコンパイル
make build-linux
```

### ローカルデータベースのセットアップ
```bash
# PostgreSQLの起動
cd local-db
docker compose up -d

# 必要な環境変数の設定
export DB_HOST=localhost
export DB_USERNAME=sbcntrapp
export DB_PASSWORD=password
export DB_NAME=sbcntrapp
export DB_CONN=1
```

## アーキテクチャ

プロジェクトはクリーンアーキテクチャに従い、以下のレイヤーで構成されています：

1. **domain/** - コアビジネスロジック
   - `model/` - ドメインエンティティとエラー型
   - `repository/` - リポジトリインターフェース

2. **usecase/** - ビジネスロジックの実装（インタラクター）

3. **handler/** - HTTPリクエストハンドラー

4. **infrastructure/** - フレームワーク固有の実装
   - ミドルウェアスタックを含むルーター設定
   - データベース接続管理

5. **interface/** - データベースインターフェースの実装

### リクエストフロー
```
HTTPリクエスト → Handler → Interactor → Repository → SQLHandler → データベース
```

### ミドルウェアスタック（実行順）
1. リクエストID生成
2. リクエストロギング（ヘルスチェックを除く）
3. AWS X-Rayトレーシング（オプション）
4. リカバリー（パニックハンドリング）

## APIエンドポイント

すべてのエンドポイントは `/v1` プレフィックス付き：

- `GET /v1/pets` - ペット一覧取得
- `POST /v1/pets/:id/like` - ペットのいいね/いいね解除
- `POST /v1/pets/:id/reservation` - ペット予約作成
- `GET /v1/notifications` - 通知一覧取得
- `POST /v1/notifications/read` - 通知を既読にする
- `GET /` と `GET /healthcheck` - ヘルスチェック
- `GET /v1/helloworld` - テストエンドポイント

## 重要な実装詳細

1. **データベース**: PostgreSQL + sqlx、安全性のため名前付きパラメータを使用
2. **トレーシング**: AWS X-Ray統合（利用不可の場合は優雅にデグレード）
3. **ロギング**: zerologによる構造化ロギング、リクエストスコープのロガー
4. **エラーハンドリング**: コード付きカスタムビジネスエラー型
5. **パフォーマンステスト**: X-Ray観測用の意図的なN+1クエリ問題とレイテンシ注入
6. **サーバーポート**: 8081（HTTP）または 443（TLS付きHTTPS）

## 一般的な開発タスク

### 新しいエンドポイントの追加
1. `domain/model/` にドメインモデルを定義
2. `domain/repository/` にリポジトリインターフェースを作成
3. `interface/` にリポジトリを実装
4. `usecase/` にインタラクターを作成
5. `handler/` にハンドラーを追加
6. `infrastructure/router.go` でルートを登録

### データベース変更
1. `db/` にマイグレーションSQLを追加
2. リポジトリインターフェースと実装を更新
3. 外部キーとよくクエリされるフィールドにはインデックスを追加することを忘れずに

### パフォーマンス問題のデバッグ
- ハンドラーとリポジトリのX-Rayサブセグメントを確認
- N+1クエリパターンをチェック（GetPetsに意図的に存在）
- コメントでマークされたランダムレイテンシ注入ポイントを監視