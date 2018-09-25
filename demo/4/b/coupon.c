#include <stdint.h>

#include "_cgo_export.h"

#include "lua.h"
#include "lauxlib.h"

#define EXECID_KEY "execid"

static uint32_t coupon_get_exec_id(lua_State *L) {
    lua_getfield(L, LUA_REGISTRYINDEX, EXECID_KEY);
    uint32_t exec_id = lua_tointeger(L, -1);
    lua_pop(L, 1);
    return exec_id;
}

static int coupon_c_set_data(lua_State *L) {
    uint32_t exec_id = coupon_get_exec_id(L);
    const char *key = luaL_checkstring(L, 1);
    const char *val = luaL_optstring(L, 2, NULL);

    // forward to go function
    coupon_go_set_data(exec_id, key, val);

    return 0;
}

static int coupon_c_get_data(lua_State *L) {
    uint32_t exec_id = coupon_get_exec_id(L);
    const char *key = luaL_checkstring(L, 1);

    // forward to go function
    // we own the returned mem, so don't forget to 'free' it
    char *val = coupon_go_get_data(exec_id, key);
    if(val == NULL) {
        lua_pushnil(L);
    } else {
        lua_pushstring(L, val);
        free((void*)(val)); // here, we not using it anymore
    }
    return 1;
}

static int coupon_c_get_age(lua_State *L) {
    uint32_t exec_id = coupon_get_exec_id(L);

    // forward to go function
    int age = coupon_go_get_age(exec_id);

    lua_pushinteger(L, age);
    return 1;
}

void coupon_bootstrap(lua_State *L, uint32_t exec_id) {
    lua_pushinteger(L, exec_id);
    lua_setfield(L, LUA_REGISTRYINDEX, EXECID_KEY);

    lua_register(L, "set_data", coupon_c_set_data);
    lua_register(L, "get_data", coupon_c_get_data);

    lua_register(L, "get_age", coupon_c_get_age);
}
