package database

// SQLHandler ...
type SQLHandler interface {
	Where(out interface{}, tableName string, clause string, args map[string]interface{}) error
	Scan(out interface{}, tableName string, order string) error
	Count(out *int, tableName string, clause string, args map[string]interface{}) error
	Create(interface{}) error
	Update(out interface{}, tableName string, clause string, args map[string]interface{}) error
}
