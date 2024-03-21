package internals

import (
	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

func list_worksheets(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	var fName string
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "file_name", &fName); err != nil {
		return nil, err
	}
	excelizeFile, err := excelize.OpenFile(fName)
	if err != nil {
		DLf(err.Error())
		return starlark.None, err
	}
	defer func() {
		if err := excelizeFile.Close(); err != nil {
			DLf(err.Error())
		}
	}()
	retArr := make([]starlark.Value, 0)
	shList := excelizeFile.GetSheetList()
	for _, shName := range shList {
		retArr = append(retArr, starlark.String(shName))
	}
	ret := starlark.NewList(retArr)
	return ret, nil
} //func list_worksheets(thread *starlark.Thread,
