package database

// SQLHandler ...
type SQLHandler interface {
	Where(interface{}, string, string, ...interface{}) error
	Scan(interface{}, string, string) error
	Count(*int, string, string, ...interface{}) error
	Create(interface{}) error
	Update(interface{}, string, string, ...interface{}) error
}
