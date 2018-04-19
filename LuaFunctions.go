package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
	"github.com/faiface/pixel/text"
)

type consoleLine struct {
	text string
	lifetime int32
}

var (
	luaRenderer *text.Text
	luaLines     = []consoleLine{{"Lua virtual machine online...", 300}}
)

func luaConsolePrint(text string) {

	luaLines = append(luaLines, consoleLine{text, 300})

	if len(luaLines) >= 32 {
		luaLines = luaLines[len(luaLines)-32:]
	}
}

func luaPrint(L *lua.LState) int {

	text := L.ToString(1)
	luaConsolePrint(text)

	return 0
}

func linkToLua(luaState *lua.LState, goFunction lua.LGFunction, goFunctionName string) {
	luaState.SetGlobal(goFunctionName, luaState.NewFunction(goFunction))
}

func executeLua(luaState *lua.LState, luaCode string) {
	err := luaState.DoString(luaCode)
	if err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFile(luaState *lua.LState, luaFile string) {
	err := luaState.DoFile(luaFile)
	if err != nil {
		fmt.Println("Lua error: " + err.Error())
	}
}

func executeLuaFunction(luaState *lua.LState, functionName string, functionArgs []lua.LValue) lua.LValue {

	err := luaState.CallByParam(lua.P{
		Fn:      luaState.GetGlobal(functionName),
		NRet:    1,
		Protect: true,
	}, functionArgs...)

	if err != nil {
		fmt.Println("Lua error: " + err.Error())
		return lua.LNil
	}

	str, ok := luaState.Get(-1).(lua.LValue)

	if ok {
		defer luaState.Pop(1)
		return str
	}

	return lua.LNil
}

func luaDisableGlobals(luaState *lua.LState) {
	executeLua(luaState, ``) //REMOVED
}