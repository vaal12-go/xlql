package internals

import (
	"fmt"

	"go.starlark.net/starlark"
)

// [x]: review if this file is needed
type StarlarkValueImplementationStub struct {
}

func (self StarlarkValueImplementationStub) Type() string {
	return "StarlarkValueImplementationStub"
}

func (self StarlarkValueImplementationStub) Freeze() {
	return
}

func (self StarlarkValueImplementationStub) Truth() starlark.Bool {
	return starlark.True
}

func (self StarlarkValueImplementationStub) Hash() (uint32, error) {
	return uint32(1), nil
}

func (self StarlarkValueImplementationStub) String() string {
	return fmt.Sprintf("StarlarkValueImplementationStub")
}

func (self StarlarkValueImplementationStub) AttrNames() []string {
	ret := []string{}
	return ret
}
