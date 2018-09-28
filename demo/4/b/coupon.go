package main

/*
#cgo CFLAGS:  -I${SRCDIR}/../_luajit/src
#cgo LDFLAGS: -L${SRCDIR}/../_luajit/src -l:libluajit.a -ldl -lm

#include <stdlib.h>
#include <stdint.h>

#include "lua.h"
#include "lauxlib.h"

#include "coupon.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func executeCoupon(s storage, uid int, couponName string, commit bool) int {
	c, err := s.getCouponForUser(uid, couponName)
	if err != nil {
		fmt.Println("DEBUG", "FAIL s.getCouponForUser", err)
		return 0
	}

	execData := &execstorageentry{
		uid:     uid,
		cdataid: c.cdataid,
		s:       s,
		data:    c.data,
	}
	execID := execstorageNew(execData)
	defer execstorageDelete(execID)

	L := C.luaL_newstate()
	defer C.lua_close(L)

	C.coupon_bootstrap(L, C.uint32_t(execID))

	codeCstr := C.CString(c.code)          // we own this mem
	defer C.free(unsafe.Pointer(codeCstr)) // so it's our responsibility to free it

	var retval C.int

	retval = C.luaL_loadstring(L, codeCstr)
	if int(retval) != 0 {
		fmt.Printf("DEBUG FAIL C.luaL_loadstring: %s\n", couponHelperLuaTostring(L, C.int(-1)))
		return 0
	}

	retval = C.lua_pcall(L, 0, 1, 0)
	if int(retval) != 0 {
		fmt.Printf("DEBUG FAIL C.lua_pcall: %s\n", couponHelperLuaTostring(L, C.int(-1)))
		return 0
	}

	result := int(C.lua_tointeger(L, -1))
	C.lua_settop(L, -2)

	if commit {
		execData.Lock()
		defer execData.Unlock()
		if err := s.saveCouponData(execData.cdataid, execData.data); err != nil {
			fmt.Println("DEBUG", "FAIL s.saveCouponData", err)
		}
	}

	return result
}

// this function is helper to get golang string type from lua stack
func couponHelperLuaTostring(state *C.lua_State, index C.int) string {
	var errmsgLen C.size_t
	errmsgPtr := C.lua_tolstring(state, index, &errmsgLen)
	if errmsgPtr == nil {
		return ""
	}
	return C.GoStringN(errmsgPtr, C.int(errmsgLen))
}

//export coupon_go_set_data
func coupon_go_set_data(execID uint32, keyCstr, valCstr *C.char) {
	execData := execstorageGet(execID)
	execData.Lock()
	defer execData.Unlock()

	key := C.GoString(keyCstr)
	if valCstr == nil {
		delete(execData.data, key)
	} else {
		val := C.GoString(valCstr)
		execData.data[key] = val
	}
}

//export coupon_go_get_data
func coupon_go_get_data(execID uint32, keyCstr *C.char) *C.char {
	execData := execstorageGet(execID)
	execData.RLock()
	defer execData.RUnlock()

	key := C.GoString(keyCstr)
	val := execData.data[key]
	if val == "" {
		return nil
	}

	return C.CString(val) // the caller own this memory
}

//export coupon_go_get_age
func coupon_go_get_age(execID uint32) C.int {
	execData := execstorageGet(execID)
	age, err := execData.s.getAgeByUserID(execData.uid)
	if err != nil {
		fmt.Println("DEBUG", "FAIL s.getAgeByUserID", err)
		return C.int(0)
	}
	return C.int(age)
}
