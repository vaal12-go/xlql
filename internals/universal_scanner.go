package internals

import (
	"time"
)

//HIGH: clean file
//HIGH: review/update logging with logging levels

// sql.Scanner interface: https://pkg.go.dev/database/sql#Scanner
// Example: https://gist.github.com/jmoiron/6979540
type UniversalScanner struct {
	internalValue any
	isNull        bool
}

// HIGH: this scanner to be used with printing and iterator for starlark
//
//	Done for iterator
//
// TODO: to be tested with time values
func (self *UniversalScanner) Scan(src interface{}) error {
	if src == nil {
		self.isNull = true
		self.internalValue = nil
		return nil
	} else {
		self.isNull = false
	}
	switch src.(type) {
	case int64:
		self.internalValue = src.(int64)
	case float64:
		self.internalValue = src.(float64)
	case bool:
		self.internalValue = src.(bool)
	case []byte:
		self.internalValue = "UniversalScanner does not implement scanning of [byte] values"
	case string:
		self.internalValue = src.(string)
	case time.Time:
		self.internalValue = src.(time.Time)
	}
	return nil
} //func (self *UniversalScanner) Scan(src interface{}) error {
