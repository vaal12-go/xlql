package internals

import (
	"go.starlark.net/starlark"
)

var MOCK_getParametersLoad_excel_sheetTESTPARAMS *loadExceSheetParams

func MOCK_getParametersLoad_excel_sheet(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	MOCK_getParametersLoad_excel_sheetTESTPARAMS, err := getParameters(args, kwargs)
	// fmt.Printf("MOCK_getParametersLoad_excel_sheetTESTPARAMS: %v\n", MOCK_getParametersLoad_excel_sheetTESTPARAMS)
	if err != nil {
		return starlark.None, err
	}
	return starlark.String(MOCK_getParametersLoad_excel_sheetTESTPARAMS.SHA256()), nil
}
