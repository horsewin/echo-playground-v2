package database

// SQLHandler ...
type SQLHandler interface {
	Where(out interface{}, tableName string, clause string, args map[string]interface{}) error
	Scan(out interface{}, tableName string, order string) error
	Count(out *int, tableName string, clause string, args map[string]interface{}) error
	Create(in map[string]interface{}, tableName string) error
	Update(in map[string]interface{}, tableName string, clause string, args map[string]interface{}) error
}
