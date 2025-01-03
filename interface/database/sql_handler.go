package database

// SQLHandler ...
type SQLHandler interface {
	Where(interface{}, string, string, ...interface{}) error
	Scan(out interface{}, tableName string, order string) error
	Count(*int, string, string, ...interface{}) error
	Create(interface{}) error
	Update(interface{}, string, string, ...interface{}) error
}
