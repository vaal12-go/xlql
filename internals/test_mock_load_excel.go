package internals

import (
	"go.starlark.net/starlark"
)

var MOCK_getParametersLoad_excel_sheetTESTPARAMS *loadExceSheetParams

func MOCK_getParametersLoad_excel_sheet(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	// fmt.Printf("Have kwargs:%v\n", kwargs)
	MOCK_getParametersLoad_excel_sheetTESTPARAMS, err := getParameters(args, kwargs)
	DLf("params: %v\n", MOCK_getParametersLoad_excel_sheetTESTPARAMS)
	if err != nil {
		return starlark.None, err
	}
	return starlark.String(MOCK_getParametersLoad_excel_sheetTESTPARAMS.String()), nil
}
