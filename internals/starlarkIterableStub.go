package internals

import (
	"fmt"

	"go.starlark.net/starlark"
)

//[x]: clean file
//HIGH: review/update logging with logging levels

type StarlarkIterableImplmentationStub struct {
	iterableCounter int64
	name            string
}

var testStub = StarlarkIterableImplmentationStub{
	iterableCounter: -1,
	name:            "testCentralStub",
}

func (self StarlarkIterableImplmentationStub) Iterate() starlark.Iterator {
	fmt.Printf("\"Iterator called\": %v\n", "Iterator called")
	ret := StarlarkIteratorImplementationStub{
		counter:        4,
		activeIterable: &self,
	}
	return ret
}

func (self StarlarkIterableImplmentationStub) String() string {
	return fmt.Sprintf("StarlarkIterableImplmentationStub:%d", self.iterableCounter)
}

func (self StarlarkIterableImplmentationStub) Type() string {
	return "StarlarkIterableImplmentationStub type"
}

func (self StarlarkIterableImplmentationStub) Freeze() {
	// panic("not implemented") //LOW: Implement Freeze() function. Low priority.
	return
}

func (self StarlarkIterableImplmentationStub) Truth() starlark.Bool {
	return starlark.Bool(true)
}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y).
// Hash may fail if the value's type is not hashable, or if the value
// contains a non-hashable value. The hash is used only by dictionaries and
// is not exposed to the Starlark program.
func (self StarlarkIterableImplmentationStub) Hash() (uint32, error) {
	return 0, nil
}

type StarlarkIteratorImplementationStub struct {
	counter        int64
	activeIterable *StarlarkIterableImplmentationStub
}

// If the iterator is exhausted, Next returns false.
// Otherwise it sets *p to the current element of the sequence,
// advances the iterator, and returns true.
func (self StarlarkIteratorImplementationStub) Next(p *starlark.Value) bool {
	if self.activeIterable.iterableCounter < 6 {
		self.activeIterable.iterableCounter++
		*p = starlark.MakeInt(int(self.activeIterable.iterableCounter))
		return true
	} else {
		return false
	}
}

func (self StarlarkIteratorImplementationStub) Done() {
	return
}
