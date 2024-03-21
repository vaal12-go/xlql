package internals

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"go.starlark.net/starlark"
)

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

// TODO: make reading of DB with callback (somesing like query_internal) for this, query and saving to xl file
func (self Query) Print(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	curs, err := NewCursorInternal(self.connected_db.db_connection, self.query_sql)
	if err != nil {
		log.Fatalf("query.print. NewCursorInternal failed with error:%v\n", err)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	colNames := curs.GetColumnNames()
	t.AppendHeader(toAnyList(colNames))
	strSlice := make([]string, len(colNames))
	for curs.Next() {
		currTblArr, err := curs.GetRow()
		if err != nil {
			log.Fatalf("query.print. GetRow failed with error:%v\n", err)
		}
		for idx := range *currTblArr {
			strSlice[idx] = getStringValueFromInterface((*currTblArr)[idx])
		} //for idx := range valSlice {
		t.AppendRow(toAnyList(strSlice))
	}
	t.AppendSeparator()
	t.SetStyle(table.StyleColoredBlackOnRedWhite)
	t.SetAllowedRowLength(220)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, WidthMin: 6, WidthMax: 10},
		{Number: 2, Align: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter},
	})
	t.SetAllowedRowLength(60)
	fmt.Println("\n*****************************************************")
	qString := self.String()
	// qString = "qwe1\nqwe2\nqwe3"
	parts := strings.Split(qString, "\n")
	for _, part := range parts {
		fmt.Printf("|  %s|\n", center(part, 49))
	}
	fmt.Println("*****************************************************")
	//Example column config:
	// t.SetColumnConfigs([]table.ColumnConfig{
	// 	{
	// 		Name:     "Institutions_name_servise_providers_name_(Eng)_",
	// 		Hidden:   false,
	// 		WidthMin: 6,
	// 		WidthMax: 12,
	// 	},
	// 	{
	// 		Name:     "Institutions_name_servise_providers_name_(Ukr)",
	// 		Hidden:   false,
	// 		WidthMin: 6,
	// 		WidthMax: 12,
	// 	},
	// 	{
	// 		Name:     "Study_team_hosptial_contact_person",
	// 		Hidden:   false,
	// 		WidthMin: 6,
	// 		WidthMax: 12,
	// 	},
	// })
	t.Render()
	return starlark.None, nil
	//[x]: review below and clean/delete
} //func (self Query) print(thread *starlark.Thread,
