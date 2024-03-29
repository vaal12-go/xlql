package internals

//DOcumentation: https://xuri.me/excelize/en/base/installation.html#read
//GitHub: https://github.com/qax-os/excelize

import (
	"log"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/xuri/excelize/v2"
)

func getColumnNamesCoords(params *loadExceSheetParams) (x_start int64, y_start int64, colNumber int64) {
	if params.skip_rows > 0 { //skip_rows are used
		return 1, params.skip_rows + 1, -1
	} else if (params.table_range_start_x > 0) ||
		(params.table_range_start_y > 0) { //table_range is used
		// fmt.Printf("\"table_range is used\": %v\n", "table_range is used")
		return params.table_range_start_x, params.table_range_start_y,
			params.table_range_end_x - params.table_range_start_x + 1
	} else { //No additional parameters should start table from  A1
		return 1, 1, -1
	} //if params.skip_rows > 0 { //skip_rows are used/
	return -1, -1, -1
} //func getColumnNamesCoords(params *loadExceSheetParams) (x_start int64, y_start int64, colNumber int64) {

func getColTypes(values_start_x int64,
	values_start_y int64, values_number int64,
	params *loadExceSheetParams) []string {
	retArr := make([]string, 0)
	for colNo := 0; colNo < int(values_number); colNo++ {
		cellCoord, err := excelize.CoordinatesToCellName(int(values_start_x)+colNo, int(values_start_y))
		if err != nil {
			log.Fatal("getColumnNamesAndTypes cannot convert cell coordinates:" + err.Error())
		}
		cellType, err := params.excelizeFile.GetCellType(params.sheet_name, cellCoord)
		if err != nil {
			log.Fatal(err)
		}
		colCell, err := params.excelizeFile.GetCellValue(
			params.sheet_name,
			cellCoord)
		if err != nil {
			log.Fatal("getColumnNamesAndTypes cannot get cell value:" + err.Error())
		}
		// DLf("colCell: %v\n", colCell)
		if cellType == 6 { //Try to parse as number if fails - parse as date
			DLf("Have a numeric format")
			t, err := dateparse.ParseAny(colCell)
			if err != nil {
				DLf("Cannot parse date")
			} else {
				DLf("Have date %s", t)
			}
			floatVal, err := strconv.ParseFloat(colCell, 64)
			if err != nil { //Can convert to Float
				DLf("Cannot convert to float")
			} else {
				DLf("Float value:%d", floatVal)
			}
		} //if cellType == 6 {
		retArr = append(retArr, getSQLTypeFromExcelType(int(cellType)))
	} //for colNo := 0; colNo < int(values_number); colNo++ {
	return retArr
} //func getColTypes(values_start_x int64,

// TODO: to rethink structure of this maybe scanning of excel and then getting names and columns types should be separated
// Also separate column scanning will require a rewrite
func getColumnNamesAndTypes(params *loadExceSheetParams) (colNames []string,
	colTypes []string, err error) {
	rows, err := params.excelizeFile.Rows(params.sheet_name)
	if err != nil {
		log.Fatal(err)
	}
	colNamesArr := make([]string, 0)
	colTypesArr := make([]string, 0)
	colNamesStartX, colNamesStartY, colNumber := getColumnNamesCoords(params)
	for rowNo := int64(1); rows.Next(); rowNo++ {
		if rowNo == colNamesStartY {
			row, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}
			if colNumber == -1 {
				colNumber = int64(len(row) - int(colNamesStartX-1))
			}
			for colNo := 0; colNo < int(colNamesStartX+colNumber); colNo++ {
				if ((colNo + 1) >= int(colNamesStartX)) &&
					(colNo+1) < int(colNamesStartX+colNumber) {
					cellCoord, err := excelize.CoordinatesToCellName(colNo+1, int(rowNo))
					if err != nil {
						log.Fatal("getColumnNamesAndTypes cannot convert cell coordinates:" + err.Error())
					}
					// fmt.Printf("cellCoord: '%v'\n", cellCoord)
					// fmt.Printf("params.sheet_name: '%v'\n", params.sheet_name)
					colCell, err := params.excelizeFile.GetCellValue(
						params.sheet_name,
						cellCoord)
					if err != nil {
						log.Fatal("getColumnNamesAndTypes cannot get cell value:" + err.Error())
					}
					colNamesArr = append(colNamesArr, colCell)
				} //if colNo == int(colNamesStartY) {
			} //for colNo := 0; colNo < int(colNumber); colNo++ {
		} //if rowNo == int(colNamesStartX) { //we are in the col row
		if rowNo == (colNamesStartY + 1) { //we are in next row, where types can be taken from
			colTypesArr = getColTypes(colNamesStartX, rowNo, colNumber, params)
			break
		}
	} //for rowNo := int64(1); rows.Next(); rowNo++ {
	rows.Close()
	//[x]: this probably is not needed - check with empty cells to the right of the 1st row of data.
	for len(colTypesArr) < int(colNumber) {
		// fmt.Printf("\"Appending type row TEXT\": %v\n", "Appending type row TEXT")
		colTypesArr = append(colTypesArr, "TEXT")
	}
	return colNamesArr, colTypesArr, nil
} //func getColumnNamesAndTypes(params *loadExceSheetParams) (colNames []string,
