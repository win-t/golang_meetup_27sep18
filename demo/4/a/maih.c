#include "_cgo_export.h"
#include "lua.h"

int luaf_c_print(lua_State *L) {
    return luaf_go_print(L);
}

void register_print(lua_State *L) {
    lua_register(L, "print", luaf_c_print);
}
