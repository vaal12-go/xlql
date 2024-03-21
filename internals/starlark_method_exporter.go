package internals

import (
	"reflect"
	"strings"

	"go.starlark.net/starlark"
)

type MethodExporter struct {
	builtInsMap map[string]*starlark.Builtin
}

func (self *MethodExporter) RegisterBuiltIns(
	object *(interface{}), methodNamesArr []string) error {
	ret_val := reflect.ValueOf(*object)
	self.builtInsMap = make(map[string]*starlark.Builtin)
	for _, method := range methodNamesArr {
		found_method := ret_val.MethodByName(method)
		if !found_method.IsValid() {
			ErrorLogger.Fatalf("ERROR: method [%s] is not implemented ", method)
		}
		builtIn := starlark.NewBuiltin(
			strings.ToLower(method), found_method.Interface().((func(thread *starlark.Thread, fn *starlark.Builtin,
				args starlark.Tuple,
				kwargs []starlark.Tuple) (starlark.Value, error))))
		builtIn.BindReceiver((*object).(starlark.Value))
		self.builtInsMap[strings.ToLower(method)] = builtIn
	}
	return nil
}

func (self MethodExporter) GetMethod(name string) (starlark.Value, error) {
	builtIn, ok := self.builtInsMap[name]
	if !ok {
		return nil, starlark.NoSuchAttrError("(self MethodExporter) GetMethod No attr name:" + name)
	} else {
		return builtIn, nil
	}
}

//TODO: add ListMethods to be used in AttrNames() functions
