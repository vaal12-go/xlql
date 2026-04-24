package internals

import (
	"go.starlark.net/starlark"
)

type StarlarkValueImplementationStub struct {
}

func (self StarlarkValueImplementationStub) Type() string {
	return "StarlarkValueImplementationStub"
}

func (self StarlarkValueImplementationStub) Freeze() {}

func (self StarlarkValueImplementationStub) Truth() starlark.Bool {
	return starlark.True
}

func (self StarlarkValueImplementationStub) Hash() (uint32, error) {
	return uint32(1), nil
}

func (self StarlarkValueImplementationStub) String() string {
	return "StarlarkValueImplementationStub"
}

func (self StarlarkValueImplementationStub) AttrNames() []string {
	ret := []string{}
	return ret
}
