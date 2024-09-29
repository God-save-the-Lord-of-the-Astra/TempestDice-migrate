package dice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func LuaReplyLoad(d *Dice) (map[string]string, error) {
	ReplyLuaCodeMap := make(map[string]string)
	files, err := filepath.Glob(d.GetExtDataDir("reply") + "/*.json")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var jsonData map[string]string
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return nil, err
		}

		for key, value := range jsonData {
			ReplyLuaCodeMap[key] = value
		}
	}

	return ReplyLuaCodeMap, nil
}

func RegisterBuiltinLuaReply(d *Dice) {
	ReplyMap, err := LuaReplyLoad(d)
	// 确保在函数结束时关闭文件
	LuaReplyExt := &ExtInfo{
		Name:       "luacommandreply", // 扩展的名称，需要用于开启和关闭指令中，写简短点
		Version:    "1.0.0",
		Brief:      "lua指令回复模块",
		Author:     "海棠",
		AutoActive: true, // 是否自动开启
		Official:   true,
		OnNotCommandReceived: func(ctx *MsgContext, msg *Message) {
			if !ctx.Dice.CustomLuaCommandConfigEnable {
				fmt.Println("1")
			} else {
				fmt.Println("0")
			}
			luaVM := lua.NewState()
			defer luaVM.Close()
			LuaVarInitWithoutArgs(luaVM, d, ctx, msg)
			LuaFuncInit(luaVM)
			if err != nil {
				ReplyToSender(ctx, msg, fmt.Sprintf("Lua 代码读取出错:\n%s", err))
				return
			} else {
				cleanText, _ := AtParse(msg.Message, "")
				cleanText = strings.TrimSpace(cleanText)
				if ReplyMap[cleanText] != "" {
					code := ReplyMap[cleanText]
					if err := luaVM.DoString(fmt.Sprintf("%s %s %s", "function main() ", code, " end")); err != nil {
						ReplyToSender(ctx, msg, fmt.Sprintf("Lua 代码执行出错:\n%s", err))
					}

					// 获取 Lua 中的 `main` 函数
					luaMain := luaVM.GetGlobal("main")

					// 调用 Lua 函数
					luaVM.Push(luaMain) // 将函数压入栈
					luaVM.Call(0, 1)    // 调用函数，0个参数，期望1个返回值

					// 获取并打印返回值
					if luaVM.GetTop() >= 1 {
						ReplyToSender(ctx, msg, fmt.Sprintf("%s%s", "Lua 代码执行成功，返回结果:\n", luaVM.ToString(-1))) // Lua栈中的最后一个元素（即返回值）
					}
				}
			}

			return
		},
		GetDescText: GetExtensionDesc,
		CmdMap:      CmdMapCls{},
	}

	d.RegisterExtension(LuaReplyExt)
}
