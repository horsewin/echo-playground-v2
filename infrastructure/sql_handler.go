package infrastructure

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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

		conn, err := sqlx.Connect("postgres", CONNECT)
		if err != nil {
			log.Fatalf("Error: No database connection established: %v", err)
		}

		// 接続成功
		fmt.Println("DB connected successfully")

		sqlHandlerInstance = &SQLHandler{Conn: conn}
	})

	return sqlHandlerInstance
}

// Where ...
func (handler *SQLHandler) Where(out interface{}, table string, whereClause string, whereArgs map[string]interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	stmt, err := handler.Conn.PrepareNamed(query)
	if err != nil {
		return err
	}

	err = stmt.Select(out, whereArgs)
	return err
}

// Scan ...
func (handler *SQLHandler) Scan(out interface{}, table string, order string) error {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", table, order)
	return handler.Conn.Select(out, query)

}

// Count ...
func (handler *SQLHandler) Count(out *int, table string, whereClause string, whereArgs map[string]interface{}) error {
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
func (handler *SQLHandler) Create(input map[string]interface{}, table string) error {
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
func (handler *SQLHandler) Update(in map[string]interface{}, table string, whereClause string, whereArgs map[string]interface{}) error {

	columns, _, values := buildNamedParameters(in)

	setClauses := make([]string, len(columns))
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf("%s = %v", col, values[col])
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ","), whereClause)
	fmt.Println(query)

	_, err := handler.Conn.NamedExec(query, whereArgs)

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
