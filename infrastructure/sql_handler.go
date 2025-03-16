package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB ...
type DB struct {
	Host     string
	Username string
	Password string
	DBName   string
	Connect  *sqlx.DB
}

// SQLHandler ... SQL handler struct
type SQLHandler struct {
	Conn *sqlx.DB
}

var (
	sqlHandlerInstance *SQLHandler
	once               sync.Once
)

const (
	dbType = "postgres"
)

// NewSQLHandler ...
func NewSQLHandler() *SQLHandler {
	once.Do(func() {
		c := utils.NewConfigDB()
		USER := c.Postgres.Username
		PASS := c.Postgres.Password
		DBNAME := c.Postgres.DBName
		PROTOCOL := "host=" + os.Getenv("DB_HOST") + " port=5432"
		CONNECT := "user=" + USER + " password=" + PASS + " " + PROTOCOL + " dbname=" + DBNAME + " sslmode=disable"

		// X-Ray対応のSQLコンテキストを作成
		db, err := xray.SQLContext(dbType, CONNECT)
		if err != nil {
			log.Fatalf("Error: No database connection established: %v", err)
		}
		conn := sqlx.NewDb(db, dbType)
		err = conn.Ping()
		if err != nil {
			db.Close()
			log.Fatalf("Error: No database connection established: %v", err)
		}

		// 接続成功
		fmt.Println("DB connected successfully")

		sqlHandlerInstance = &SQLHandler{Conn: conn}
	})

	return sqlHandlerInstance
}

// Where ...
func (handler *SQLHandler) Where(ctx context.Context, out interface{}, table string, whereClause string, whereArgs map[string]interface{}) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Where")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Where")
		return handler.whereWithoutXRay(out, table, whereClause, whereArgs)
	}
	defer seg.Close(nil)

	query := fmt.Sprintf("SELECT * FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("args", whereArgs); err != nil {
		utils.LogError("Failed to add args metadata: %v", err)
	}

	stmt, err := handler.Conn.PrepareNamed(query)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
		return err
	}

	err = stmt.SelectContext(subCtx, out, whereArgs)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// whereWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) whereWithoutXRay(out interface{}, table string, whereClause string, whereArgs map[string]interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	stmt, err := handler.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.Select(out, whereArgs)
}

// Scan ...
func (handler *SQLHandler) Scan(ctx context.Context, out interface{}, table string, order string) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Scan")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Scan")
		return handler.scanWithoutXRay(out, table, order)
	}
	defer seg.Close(nil)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", table, order)

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}

	err := handler.Conn.SelectContext(subCtx, out, query)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// scanWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) scanWithoutXRay(out interface{}, table string, order string) error {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", table, order)
	return handler.Conn.Select(out, query)
}

// Count ...
func (handler *SQLHandler) Count(ctx context.Context, out *int, table string, whereClause string, whereArgs map[string]interface{}) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Count")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Count")
		return handler.countWithoutXRay(out, table, whereClause, whereArgs)
	}
	defer seg.Close(nil)

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("args", whereArgs); err != nil {
		utils.LogError("Failed to add args metadata: %v", err)
	}

	var count int
	stmt, err := handler.Conn.PrepareNamed(query)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
		return err
	}

	err = stmt.GetContext(subCtx, &count, whereArgs)
	*out = count
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// countWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) countWithoutXRay(out *int, table string, whereClause string, whereArgs map[string]interface{}) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	var count int
	stmt, err := handler.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}

	err = stmt.Get(&count, whereArgs)
	*out = count
	return err
}

// Create ...
func (handler *SQLHandler) Create(ctx context.Context, input map[string]interface{}, table string) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Create")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Create")
		return handler.createWithoutXRay(input, table)
	}
	defer seg.Close(nil)

	// カラム名とプレースホルダーを構築
	columns := make([]string, 0)
	placeholders := make([]string, 0)

	// inputのキーと値をそれぞれ列と値に追加
	for key := range input {
		// IDが設定されていても無視する
		if key == "id" {
			continue
		}

		columns = append(columns, key)
		placeholders = append(placeholders, fmt.Sprintf(":%s", key)) // プレースホルダーを使う
	}

	// クエリを構築
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Join(placeholders, ","))

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("input", input); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	_, err := handler.Conn.NamedExecContext(subCtx, query, input)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// createWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) createWithoutXRay(input map[string]interface{}, table string) error {
	// カラム名とプレースホルダーを構築
	columns := make([]string, 0)
	placeholders := make([]string, 0)

	// inputのキーと値をそれぞれ列と値に追加
	for key := range input {
		// IDが設定されていても無視する
		if key == "id" {
			continue
		}

		columns = append(columns, key)
		placeholders = append(placeholders, fmt.Sprintf(":%s", key)) // プレースホルダーを使う
	}

	// クエリを構築
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Join(placeholders, ","))

	fmt.Println(query)
	_, err := handler.Conn.NamedExec(query, input)
	return err
}

// Update ...
func (handler *SQLHandler) Update(ctx context.Context, in map[string]interface{}, table string, whereClause string) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Update")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Update")
		return handler.updateWithoutXRay(in, table, whereClause)
	}
	defer seg.Close(nil)

	columns, placeholders, _ := buildNamedParameters(in)

	setClauses := make([]string, len(columns))
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf("%s = %s", col, placeholders[i])
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ","), whereClause)

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("input", in); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	fmt.Println(query)

	_, err := handler.Conn.NamedExecContext(subCtx, query, in)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}

	return err
}

// updateWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) updateWithoutXRay(in map[string]interface{}, table string, whereClause string) error {
	columns, placeholders, _ := buildNamedParameters(in)

	setClauses := make([]string, len(columns))
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf("%s = %s", col, placeholders[i])
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ","), whereClause)

	fmt.Println(query)

	_, err := handler.Conn.NamedExec(query, in)
	return err
}

// Delete ...
func (handler *SQLHandler) Delete(ctx context.Context, in map[string]interface{}, table string) error {
	// X-Rayサブセグメントを作成
	subCtx, seg := xray.BeginSubsegment(ctx, "SQLHandler.Delete")
	if seg == nil {
		// セグメントが作成できない場合はログに記録して処理を続行
		utils.LogError("Failed to begin subsegment: SQLHandler.Delete")
		return handler.deleteWithoutXRay(in, table)
	}
	defer seg.Close(nil)

	columns, _, values := buildNamedParameters(in)

	whereClauses := make([]string, len(columns))
	for i, col := range columns {
		whereClauses[i] = fmt.Sprintf("%s = %v", col, values[col])
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(whereClauses, ","))

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}
	if err := seg.AddMetadata("input", in); err != nil {
		utils.LogError("Failed to add input metadata: %v", err)
	}

	fmt.Println(query)

	_, err := handler.Conn.NamedExecContext(subCtx, query, values)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}

	return err
}

// deleteWithoutXRay はX-Rayなしでクエリを実行するためのヘルパーメソッド
func (handler *SQLHandler) deleteWithoutXRay(in map[string]interface{}, table string) error {
	columns, _, values := buildNamedParameters(in)

	whereClauses := make([]string, len(columns))
	for i, col := range columns {
		whereClauses[i] = fmt.Sprintf("%s = %v", col, values[col])
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, strings.Join(whereClauses, ","))

	fmt.Println(query)

	_, err := handler.Conn.NamedExec(query, values)
	return err
}

func buildNamedParameters(input map[string]interface{}) (columns []string, placeholderNames []string, values map[string]interface{}) {
	columns = []string{}
	placeholderNames = []string{}
	values = make(map[string]interface{})

	for key, value := range input {
		if value != nil {
			columns = append(columns, key)
			placeholderNames = append(placeholderNames, ":"+key)
			values[key] = value
		}
	}

	return
}
