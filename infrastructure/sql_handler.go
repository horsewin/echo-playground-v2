package infrastructure

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/horsewin/echo-playground-v2/utils"
	"github.com/jmoiron/sqlx"
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

// NewSQLHandler ...
func NewSQLHandler() *SQLHandler {
	c := utils.NewConfigDB()
	USER := c.MySQL.Username
	PASS := c.MySQL.Password
	PROTOCOL := c.MySQL.Protocol
	DBNAME := c.MySQL.DBName
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true&loc=Asia%2FTokyo"

	conn, err := sqlx.Connect("mysql", CONNECT)
	if err != nil {
		log.Fatalf("Error: No database connection established: %v", err)
	}

	sqlHandler := new(SQLHandler)
	sqlHandler.Conn = conn
	return sqlHandler
}

// Where ...
func (handler *SQLHandler) Where(out interface{}, table string, query string, args ...interface{}) error {
	if query == "" {
		query = fmt.Sprintf("SELECT * FROM %s", table)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s", table, query)
	}
	return handler.Conn.Select(out, query, args...)
}

// Scan ...
func (handler *SQLHandler) Scan(out interface{}, table string, order string) error {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", table, order)
	return handler.Conn.Select(out, query)

}

// Count ...
func (handler *SQLHandler) Count(out *int, table string, whereClause string, whereArgs ...interface{}) error {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if whereClause != "" {
		query += fmt.Sprintf(" WHERE %s", whereClause)
	}

	var count int
	err := handler.Conn.Get(&count, query, whereArgs...)
	*out = count
	return err
}

// Create ...
func (handler *SQLHandler) Create(input interface{}) error {
	table := reflect.TypeOf(input).Elem().Name() // ポインタをdereference
	columns, values, args := buildNamedParameters(input)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Join(values, ","))

	_, err := handler.Conn.NamedExec(query, input)
	return err
}

// Update ...
func (handler *SQLHandler) Update(input interface{}, table string, whereClause string, whereArgs ...interface{}) error {

	columns, values, _ := buildNamedParameters(input)

	setClauses := make([]string, len(columns))
	for i, col := range columns {
		setClauses[i] = fmt.Sprintf("%s = :%s", col, col)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(setClauses, ","), whereClause)

	_, err := handler.Conn.NamedExec(query, input)

	return err
}

// 変更点：map[string]interface{} を返すように変更
func buildNamedParameters(input interface{}) (columns []string, values []string, args map[string]interface{}) {
	columns = []string{}
	values = []string{}
	args = make(map[string]interface{})

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		if !v.Field(i).IsZero() {
			columns = append(columns, tag)
			values = append(values, ":"+tag)
			args[tag] = v.Field(i).Interface() // 値をマップに追加
		}
	}

	return
}
