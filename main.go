package main

//Example: https://otm.github.io/2015/07/embedding-lua-in-go/
//Package: https://github.com/yuin/gopher-lua
//Documentation: https://pkg.go.dev/github.com/yuin/gopher-lua#LState

import (
	lua "github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("Hello World")`); err != nil {
		panic(err)
	}
}
