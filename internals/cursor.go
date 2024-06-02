package internals

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.starlark.net/starlark"
)

type Cursor struct {
	query               *Query
	colNames            []string
	rows                *sql.Rows
	cursor_ref          *Cursor //This is hack to make Next method be able to change fields of this struct
	cursor_done_builtin *starlark.Builtin
	isClosed            bool
	cursor_internal     *cursor_internal
}

func (self *Cursor) cursor_done(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	self.Done()
	return starlark.None, nil
}

func NewCursor(query *Query) *Cursor {
	ret := Cursor{
		query:    query,
		isClosed: false,
	}
	var err error
	ret.cursor_internal, err = NewCursorInternal(
		query.connected_db.db_connection,
		query.query_sql)
	if err != nil {
		log.Fatalf("NewCursor failed creating cursor_internal:%v", err)
	}
	ret.cursor_done_builtin = starlark.NewBuiltin("cursor_done", ret.cursor_done)
	ret.cursor_done_builtin.BindReceiver(ret)
	ret.colNames = ret.cursor_internal.GetColumnNames()
	if err != nil {
		log.Fatal(err)
	}
	ret.cursor_ref = &ret
	return &ret
} //func NewCursor(query *Query) *Cursor {

func (self Cursor) Iterate() starlark.Iterator {
	return self
}

func (self Cursor) Next(p *starlark.Value) bool {
	if self.cursor_internal.Next() {
		valSlice, _ := self.cursor_internal.GetRow()
		retArr := make([]starlark.Value, 0)
		for _, val := range *valSlice {
			if val == nil {
				retArr = append(retArr, starlark.None)
				continue
			}
			t1 := reflect.TypeOf(val)
			switch t := val.(type) {
			case string:
				retArr = append(retArr,
					starlark.String((val.(string))))
			case float64:
				retArr = append(retArr,
					starlark.Float((val.(float64))))
			case int64:
				retArr = append(retArr,
					starlark.MakeInt64(val.(int64)))
			case time.Time:
				retArr = append(retArr,
					starlark.String((val.(time.Time).String())))

			default:
				fmt.Printf("unknown type:%s     %v\n", t1, val)
				log.Fatalf("cursor.Next(). Handling of this type is not implemented:%v", t)
			}
		} ////for _, val := range *valSlice {
		*p = starlark.NewList(retArr)
		return true
	} else { //if self.cursor_internal.Next() {
		return false
	}
} //func (self Cursor) Next(p *starlark.Value) bool {

func (self Cursor) Done() {
	if !self.isClosed {
		self.cursor_internal.Close()
		self.cursor_ref.isClosed = true
	} else {
		DLf("Second time closing cursor - will do nothing")
	}
	return
} //func (self Cursor) Done() {

func (self Cursor) Attr(name string) (starlark.Value, error) {
	switch name {
	case "cursor_done":
		return self.cursor_done_builtin, nil
	}
	return starlark.None, nil
}

func (self Cursor) AttrNames() []string {
	return []string{"cursor_done"}
}

// TODO: implement starlark object methods
func (self Cursor) Type() string {
	return "Cursor"
}

func (self Cursor) Freeze() {
	return
}

func (self Cursor) Truth() starlark.Bool {
	return starlark.True
}

func (self Cursor) Hash() (uint32, error) {
	return uint32(1), nil
}

func (self Cursor) String() string {
	return fmt.Sprintf("Cursor")
}
