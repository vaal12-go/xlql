package internals

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

//[x]: clean file

func calcExcelLoadableRange(params *loadExceSheetParams) error {
	if params.table_range_start != "" { //This is table range situation
		if params.skip_rows > 0 {
			return fmt.Errorf("calcExcelLoadableRange. table_range cannot be used with skip_rows in load_excel_sheet")
		}
		var colNumber int64
		x_start, y_start, err := excelize.CellNameToCoordinates(params.table_range_start)
		if err != nil {
			return err
		}
		_, y_end, err := excelize.CellNameToCoordinates(params.table_range_end)
		if err != nil {
			return err
		}
		if y_end > y_start { //There is data rows
			params.calc_data_start_row = int64(y_start) + 1
			params.calc_data_end_row = int64(y_end)
			params.calc_data_start_col = int64(x_start)
		}
		//TODO: add checks that x\y start are less or equal than x\y end

		params.calc_header_start_col, params.calc_header_start_row, colNumber = getColumnNamesCoords(params)
		params.calc_header_end_col = params.calc_header_start_col + colNumber - 1
	} else { //Simple table situation
		//TODO: add support for skip_rows
		if params.skip_rows > 0 {
			params.calc_header_start_row = params.skip_rows + 1
		} else {
			params.calc_header_start_row = 1
		}
		params.calc_header_start_col = 1
		rows, err := params.excelizeFile.Rows(params.sheet_name)
		if err != nil {
			return err
		}
		rows.Next()
		cells, err := rows.Columns()
		params.calc_header_end_col = params.calc_header_start_col + int64(len(cells)) - 1
		data_row_no := 0
		for rows.Next() {
			if data_row_no == 0 { //This is first separate data row
				params.calc_data_start_row = params.calc_header_start_row + 1
				params.calc_data_start_col = params.calc_header_start_col
				params.calc_data_end_row = 0
			}
			params.calc_data_end_row++
			data_row_no++
		}
	} //} else { //Simple table situation
	return nil
} //func calcExcelLoadableRange(params *loadExceSheetParams) error {
