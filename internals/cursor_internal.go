package internals

import (
	"context"
	"database/sql"
)

type cursor_internal struct {
	dbConnection *sql.DB
	query_string string
	rows         *sql.Rows
	colNames     []string
	colTypes     ([](*sql.ColumnType))
}

func NewCursorInternal(conn *sql.DB, query string) (*cursor_internal, error) {
	ret := cursor_internal{
		dbConnection: conn,
		query_string: query,
	}
	var err error
	tx, err := ret.dbConnection.BeginTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  true,
	})
	ret.rows, err = ret.dbConnection.Query(ret.query_string)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	ret.colNames, err = ret.rows.Columns()
	if err != nil {
		return nil, err
	}
	ret.colTypes, err = ret.rows.ColumnTypes()
	return &ret, nil
} //func NewCursorInternal(conn *sql.DB, query string) (*cursor_internal, error) {

func (self *cursor_internal) GetColumnNames() []string {
	//If returned as is, array can be modified externally need to return copy
	ret := make([]string, len(self.colNames))
	copy(ret, self.colNames)
	return ret
}

func (self *cursor_internal) GetColumnTypes() []*sql.ColumnType {
	//If returned as is, array can be modified externally need to return copy
	ret := make([]*sql.ColumnType, len(self.colTypes))
	copy(ret, self.colTypes)
	return ret
}

func (self *cursor_internal) GetRow() (*([](interface{})), error) {
	valSlice := make([]interface{}, len(self.colNames))
	colTypes, err := self.rows.ColumnTypes()
	for idx, _ := range colTypes { //Creating type specific receiver array for values
		valSlice[idx] = new(UniversalScanner)
	} //for idx, colType := range colTypes {//Creating type specific receiver array for values
	err = self.rows.Scan(valSlice...)
	if err != nil {
		return nil, err
	}
	ret := make([](interface{}), len(self.colNames))
	for idx, val := range valSlice { //Converting received values to actual values (handing NULLs)
		ret[idx] = (val.(*UniversalScanner)).internalValue
	}
	return &ret, nil
	// return self.readDBRow()
}

// [x]: consider merging to GetRow - now function seems to be very simplified
// Read columns to slice: https://stackoverflow.com/questions/14477941/read-select-columns-into-string-in-go
// func (self *cursor_internal) readDBRow() (*([]interface{}), error) {

// } //func readDBRow(rows *sql.Rows, numCols int) *([]interface{}), error {

func (self *cursor_internal) Next() bool {
	return self.rows.Next()
}

func (self *cursor_internal) Close() {
	self.rows.Close()
}
