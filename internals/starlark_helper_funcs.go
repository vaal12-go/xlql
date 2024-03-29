package internals

import (
	"time"

	"go.starlark.net/starlark"
)

func get_datetime_formatted(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	dateFormat := "2006-01-02[15.04.05]"
	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"format?", &dateFormat); err != nil {
		return starlark.None, err
	}

	tm_str := time.Now().Format(dateFormat)
	return starlark.String(tm_str), nil
}
