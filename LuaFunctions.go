package main

import (
	"github.com/yuin/gopher-lua"
	"fmt"
	"github.com/faiface/pixel"
)

func luaPrint(L *lua.LState) int {
	text := L.ToString(1)
	luaRenderer.Dot = pixel.V(screenWidth*0.75, screenHeight-22 - float64(luaLine)*22)
	luaRenderer.WriteString(text + "\n")
	luaLine++
	return 0
}


func luaClear(L *lua.LState) int {
	luaRenderer.Clear()
	luaLine = 0
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