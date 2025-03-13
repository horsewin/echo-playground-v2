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

// NewSQLHandler ...
func NewSQLHandler() *SQLHandler {
	once.Do(func() {
		c := utils.NewConfigDB()
		USER := c.Postgres.Username
		PASS := c.Postgres.Password
		DBNAME := c.Postgres.DBName
		PROTOCOL := "host=" + os.Getenv("DB_HOST") + " port=5432"
		CONNECT := "user=" + USER + " password=" + PASS + " " + PROTOCOL + " dbname=" + DBNAME + " sslmode=disable"

		db, err := xray.SQLContext("postgres", CONNECT)
		if err != nil {
			log.Fatalf("Error: No database connection established: %v", err)
		}
		conn := sqlx.NewDb(db, "postgres")
		err = conn.Ping()
		if err != nil {
			db.Close()
			log.Fatalf("Error: No database connection established: %v", err)
			os.Exit(1)
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
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Where")
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

	err = stmt.Select(out, whereArgs)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// Scan ...
func (handler *SQLHandler) Scan(ctx context.Context, out interface{}, table string, order string) error {
	// X-Rayサブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Scan")
	defer seg.Close(nil)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", table, order)

	// クエリをメタデータとして追加
	if err := seg.AddMetadata("query", query); err != nil {
		utils.LogError("Failed to add query metadata: %v", err)
	}

	err := handler.Conn.Select(out, query)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// Count ...
func (handler *SQLHandler) Count(ctx context.Context, out *int, table string, whereClause string, whereArgs map[string]interface{}) error {
	// X-Rayサブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Count")
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

	err = stmt.Get(&count, whereArgs)
	*out = count
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// Create ...
func (handler *SQLHandler) Create(ctx context.Context, input map[string]interface{}, table string) error {
	// X-Rayサブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Create")
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

	fmt.Println(query)
	_, err := handler.Conn.NamedExec(query, input)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}
	return err
}

// Update ...
func (handler *SQLHandler) Update(ctx context.Context, in map[string]interface{}, table string, whereClause string) error {
	// X-Rayサブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Update")
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

	_, err := handler.Conn.NamedExec(query, in)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}

	return err
}

// Delete ...
func (handler *SQLHandler) Delete(ctx context.Context, in map[string]interface{}, table string) error {
	// X-Rayサブセグメントを作成
	_, seg := xray.BeginSubsegment(ctx, "SQLHandler.Delete")
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

	_, err := handler.Conn.NamedExec(query, values)
	if err != nil {
		if addErr := seg.AddError(err); addErr != nil {
			utils.LogError("Failed to add error to segment: %v", addErr)
		}
	}

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
