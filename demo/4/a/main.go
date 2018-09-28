package main

/*
#cgo CFLAGS:  -I${SRCDIR}/../_luajit/src
#cgo LDFLAGS: -L${SRCDIR}/../_luajit/src -l:libluajit.a -ldl -lm

#include <stdlib.h>

#include "lua.h"
#include "lauxlib.h"

#include "main.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	L := C.luaL_newstate()
	defer C.lua_close(L)

	C.register_print(L)

	simpleLuaProgram := `
		function fact(n, acc)
			acc = acc or 1
			if n <= 1 then
				return acc
			else
				return fact(n - 1, n * acc)
			end
		end

		print(fact(6))
		print("Hello from lua")
	`
	simpleLuaProgramPtr := C.CString(simpleLuaProgram) // we own this mem
	defer C.free(unsafe.Pointer(simpleLuaProgramPtr))  // so it's our responsibility to free it

	var retval C.int

	retval = C.luaL_loadstring(L, simpleLuaProgramPtr)
	if int(retval) != 0 {
		fmt.Printf("Cannot load program: %s\n", go_lua_tostring(L, C.int(-1)))
	}

	retval = C.lua_pcall(L, 0, 0, 0)
	if int(retval) != 0 {
		fmt.Printf("Cannot execute program: %s\n", go_lua_tostring(L, C.int(-1)))
	}
}

func go_lua_tostring(state *C.lua_State, index C.int) string {
	var errmsgLen C.size_t
	errmsgPtr := C.lua_tolstring(state, index, &errmsgLen)
	if errmsgPtr == nil {
		return ""
	}
	return C.GoStringN(errmsgPtr, C.int(errmsgLen))
}

//export luaf_go_print
func luaf_go_print(L *C.lua_State) C.int {
	t := C.lua_type(L, 1)
	switch t {
	case C.LUA_TNUMBER:
		n := C.lua_tonumber(L, 1)
		fmt.Printf("%v\n", n)
	case C.LUA_TSTRING:
		s := go_lua_tostring(L, 1)
		fmt.Println(s)
	default:
		fmt.Println("Unknown Type")
	}
	return 0
}
