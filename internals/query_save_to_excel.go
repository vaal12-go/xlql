package internals

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

// TODO: check how null values from table appear in excel
// TODO: save_to_excel: LOW: https://pkg.go.dev/github.com/xuri/excelize/v2#File.SetColWidth - use SetColWidth to set width of columns which are very large (e.g. >10 symbols)
// use SetColWidth to set width of columns (or titles) which are very large (e.g. >10 symbols)

func (self Query) Save_to_excel(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	fName, sheetName, tableStyle := "", "Sheet1", ""
	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"file_name", &fName,
		"sheet_name?", &sheetName,
		"table_style?", &tableStyle); err != nil {
		fmt.Printf("\"FIle name issue:\": %v\n %s\n", "FIle name issue:", err)
		return nil, err
	}

	_, err := os.Stat(fName)
	var excelizeFile *excelize.File
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			excelizeFile = excelize.NewFile()
		} else {
			//Some other error - will panic
			fmt.Printf("Some error checking file [%s]exists:%v\n",
				fName, err)
			return starlark.None, err
		}

	} else {
		excelizeFile, err = excelize.OpenFile(fName)
		if err != nil {
			return starlark.None, err
		}
	}
	defer func() {
		if err := excelizeFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	sheetIdx, err := excelizeFile.NewSheet(sheetName)
	if err != nil {
		ErrorLogger.Printf("save_to_excel. Error creating new worksheet:%s", err)
		return starlark.None, err
	}
	q, err := NewCursorInternal(
		self.connected_db.db_connection,
		self.query_sql)
	defer q.Close()
	if err != nil {
		ErrorLogger.Printf("NewCursor creation error: %v\n", err)
		log.Fatal(err)
	}
	colArr := q.GetColumnNames()
	typesArr := q.GetColumnTypes()
	// fmt.Printf("typesArr: %v\n", typesArr)
	// for _, tp := range typesArr {
	// 	fmt.Printf("tp: %v\n", tp)
	// }
	firstCellAddr := "A1"
	lastCellAddr := "A1"
	for idx, colName := range colArr { //Adding column names row
		cellAddrStr, err := excelize.CoordinatesToCellName(idx+1, 1)
		if err != nil {
			WarningLogger.Printf(
				"save_to_excel. Error adding column names:%v\n", err)
			return starlark.None, err
		}
		excelizeFile.SetCellValue(sheetName,
			cellAddrStr,
			colName)
		lastCellAddr = cellAddrStr
	} //for idx, colName := range colArr {//Adding column names row
	i := 0
	for currRow := 2; q.Next(); currRow++ {
		currCol := 1
		valSlice, err := q.GetRow()
		// fmt.Printf("valSlice: %v\n", valSlice)
		if err != nil {
			log.Fatal(err)
		}
		for idx := range *valSlice {
			cellAddrStr, err := excelize.CoordinatesToCellName(currCol, currRow)
			if err != nil {
				WarningLogger.Printf(
					"save_to_excel. Error populating cell values:%v\n", err)
				return starlark.None, err
			}
			//LOW: consider to add all_as_text parameter to save_to_excel to avoid type conversion
			val := ((*valSlice)[idx])
			// fmt.Printf("\nval: %v\n", val)
			// fmt.Printf("typesArr[idx].DatabaseTypeName(): %v\n", typesArr[idx].DatabaseTypeName())
			excelizeFile.SetCellValue(sheetName,
				cellAddrStr,
				(val))

			if typesArr[idx].DatabaseTypeName() == "DATE" {
				NUM_FORMAT := "d-mmm-yy h:mm"
				style, err := excelizeFile.NewStyle(&excelize.Style{CustomNumFmt: &NUM_FORMAT})

				if err != nil {
					fmt.Println(err)
				}
				err = excelizeFile.SetCellStyle(sheetName, cellAddrStr, cellAddrStr, style)
			}
			lastCellAddr = cellAddrStr
			currCol++
		} //for idx := range *valSlice {
		i++
	} //for rows.Next() {
	tblRange := fmt.Sprintf("%s:%s", firstCellAddr, lastCellAddr)
	disable := false
	enable := true
	//LOW: Add parameter to call xl_table_style with styles from https://xuri.me/excelize/en/utils.html#AddTable
	tbl_name := "table_" + sheetName
	for {
		//TODO: add ShowRowStripes, ShowColumnStripes, ShowFirst Column and Show Last column as parameters
		//TODO: see how cell level formatting can be applied: https://xuri.me/excelize/en/cell.html#SetCellStyle   See Example #5
		err = excelizeFile.AddTable(sheetName,
			&excelize.Table{
				Range: tblRange,
				Name:  tbl_name,
				// StyleName: "TableStyleMedium2",
				StyleName:         tableStyle,
				ShowFirstColumn:   true,
				ShowLastColumn:    true,
				ShowRowStripes:    &enable,
				ShowColumnStripes: disable,
			})
		if err != nil {
			if strings.Contains(err.Error(), "the same name table already exists") {
				tbl_name = tbl_name + "0"
				// fmt.Printf("tbl_name: %v\n", tbl_name)
			} else {
				return starlark.None, fmt.Errorf("Error setting table in Excel:%v", err)
			}
		} else {
			break
		}
	}

	excelizeFile.SetActiveSheet(sheetIdx)

	currTimeStr := time.Now().Format(time.RFC1123)

	err = excelizeFile.AddComment(sheetName, excelize.Comment{
		Cell:   "A1",
		Author: "Excelize",
		Paragraph: []excelize.RichTextRun{
			{Text: "XLQL: ", Font: &excelize.Font{Bold: true}},
			{Text: "This worksheet was generated by XLQL (https://github.com/vaal12-go/xlql/) \n"},
			{Text: "On:" + currTimeStr + "\n", Font: &excelize.Font{Italic: true}},
			{Text: "From a query:"},
			{Text: self.query_sql + "\n", Font: &excelize.Font{Bold: true}},
		},
	})
	if err != nil {
		fmt.Printf("Error adding metadata comment: %v\n", err)
		return starlark.None, err
	}

	//LOW: add parameter (for user to use) remove_initial_worksheet to remove "Sheet1" from the created file
	if err := excelizeFile.SaveAs(fName); err != nil {
		ErrorLogger.Printf("save_to_excel. Error saving file: %v\n", err)
		return starlark.None, err
	}
	return starlark.None, nil
} //func (self Query) save_to_excel(thread *starlark.Thread,
