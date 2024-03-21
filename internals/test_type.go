package internals

import (
	"fmt"

	"go.starlark.net/starlark"
)

//REMOVE THIS FILE TO SEPARATE BRANCH

type testType struct {
	memberFunc *starlark.Builtin
	id         int
}

// https://pkg.go.dev/go.starlark.net/starlark#Value
// Starlark string values are quoted as if by Python's repr.
func (self testType) String() string {
	// panic("not implemented") // TODO: Implement String of test_type
	return "String of testType"
}

// Type returns a short string describing the value's type.
func (self testType) Type() string {
	return "testType"
}

func (self testType) Freeze() {
	return
}

// Truth returns the truth value of an object.
func (self testType) Truth() starlark.Bool {
	return starlark.True
}

func (self testType) Hash() (uint32, error) {
	return 1, nil
}

// See comment in value.go file
func (self testType) Attr(name string) (starlark.Value, error) {
	switch name {
	case "sayhello":
		return self.memberFunc, nil
	case "id":
		return starlark.MakeInt(self.id), nil
	default:
		return nil, starlark.NoSuchAttrError("No attr name:" + name)
	}

}

func (self testType) AttrNames() []string {
	ret := []string{"sayhello"}
	return ret
}

func (self testType) Iterate() starlark.Iterator {
	newIter := testIterator{
		currIteration: 0,
	}
	fmt.Println("Iterate called")
	return newIter
}

type testIterator struct {
	currIteration int
}

var currIteration int = 0

// If the iterator is exhausted, Next returns false.
// Otherwise it sets *p to the current element of the sequence,
// advances the iterator, and returns true.
func (self testIterator) Next(p *starlark.Value) bool {
	// fmt.Println("Next called")
	if currIteration < 100 {
		*p = starlark.MakeInt(currIteration)
		currIteration++
		return true
	} else {
		return false
	}

}

func (self testIterator) Done() {
	self.currIteration = -1
}
