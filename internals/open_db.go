package internals

import (
	"go.starlark.net/starlark"
)

// [x]: clean file
//[x]: review/update logging with logging levels

func open_db(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	var fName string
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "file_name", &fName); err != nil {
		return nil, err
	}
	ret := NewDatabase(fName)
	return ret, nil
}
