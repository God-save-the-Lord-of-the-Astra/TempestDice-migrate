package dice

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	lua "github.com/yuin/gopher-lua"
)

var watcher *fsnotify.Watcher

func LuaReplyLoad(d *Dice) {
	ReplyLuaCodeMap := make(map[string]string)
	files, _ := filepath.Glob(d.GetExtDataDir("reply") + "/*.json")
	defaultLuaReplyExist := false
	var filename string
	for _, file := range files {
		if runtime.GOOS == "windows" {
			filename = "\\luareply.json"
		} else {
			filename = "/luareply.json"
		}
		if file == d.GetExtDataDir("reply")+filename {
			defaultLuaReplyExist = true
		}
		data, _ := os.ReadFile(file)

		var jsonData map[string]string
		json.Unmarshal(data, &jsonData)

		for key, value := range jsonData {
			ReplyLuaCodeMap[key] = value
		}
	}
	if !defaultLuaReplyExist {
		d.Logger.Infof("没有找到默认的luareply.json，请检查 %s 目录下是否存在", d.GetExtDataDir("reply"))
	}

	d.CustomLuaReplyMap = ReplyLuaCodeMap
}

func RegisterBuiltinLuaReply(d *Dice) {
	LuaReplyLoad(d)

	LuaReplyExt := &ExtInfo{
		Name:       "luareply",
		Version:    "1.0.0",
		Brief:      "lua指令回复模块",
		Author:     "海棠",
		AutoActive: true,
		Official:   true,
		OnMessageReceived: func(ctx *MsgContext, msg *Message) {
			luaInitStartTime := time.Now().UnixMicro()
			if !ctx.Dice.CustomReplyConfigEnable {
				return
			}
			ReplyMap := d.CustomLuaReplyMap
			luaVM := lua.NewState()
			defer luaVM.Close()
			LuaVarInitWithoutArgs(luaVM, d, ctx, msg)
			LuaFuncInit(luaVM)
			cleanText, _ := AtParse(msg.Message, "")
			cleanText = strings.TrimSpace(cleanText)

			var matchedCode string
			for pattern, code := range ReplyMap {
				matched, err := regexp.MatchString(pattern, cleanText)
				if err != nil {
					d.Logger.Error(fmt.Sprintf("正则表达式编译错误: %s", err))
					continue
				}
				if matched {
					matchedCode = code
					break
				}
			}

			if matchedCode != "" {
				code := matchedCode
				if err := luaVM.DoString(fmt.Sprintf("%s %s %s", "function main() ", code, " end")); err != nil {
					ReplyToSender(ctx, msg, fmt.Sprintf("Lua 代码执行出错:\n%s", err))
				}

				luaMain := luaVM.GetGlobal("main")
				luaVM.Push(luaMain)
				luaVM.Call(0, 1)

				if luaVM.GetTop() >= 1 {
					luaInitEndTime := time.Now().UnixMicro()
					ReplyToSender(ctx, msg, luaVM.ToString(-1))
					d.Logger.Info(fmt.Sprintf("%s%d%s", "[回复调试] lua库初始化结束，耗时: ", (luaInitEndTime-luaInitStartTime)/1000, "毫秒\n"))
				}
			}
		},
		GetDescText: GetExtensionDesc,
		CmdMap:      CmdMapCls{},
	}

	d.RegisterExtension(LuaReplyExt)
}
