package dice

import (
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"os/exec"
	"strings"
	"time"

	"github.com/samber/lo"

	ds "github.com/sealdice/dicescript"
)

var emokloreAttrParent = map[string][]string{
	"检索":   {"知力"},
	"洞察":   {"知力"},
	"识路":   {"灵巧", "五感"},
	"直觉":   {"精神", "运势"},
	"鉴定":   {"五感", "知力"},
	"观察":   {"五感"},
	"聆听":   {"五感"},
	"鉴毒":   {"五感"},
	"危机察觉": {"五感", "运势"},
	"灵感":   {"精神", "运势"},
	"社交术":  {"社会"},
	"辩论":   {"知力"},
	"心理":   {"精神", "知力"},
	"魅惑":   {"魅力"},
	"专业知识": {"知力"},
	"万事通":  {"五感", "社会"},
	"业界":   {"社会", "魅力"},
	"速度":   {"身体"},
	"力量":   {"身体"},
	"特技动作": {"身体", "灵巧"},
	"潜泳":   {"身体"},
	"武术":   {"身体"},
	"奥义":   {"身体", "精神", "灵巧"},
	"射击":   {"灵巧", "五感"},
	"耐久":   {"身体"},
	"毅力":   {"精神"},
	"医术":   {"灵巧", "知力"},
	"技巧":   {"灵巧"},
	"艺术":   {"灵巧", "精神", "五感"},
	"操纵":   {"灵巧", "五感", "知力"},
	"暗号":   {"知力"},
	"电脑":   {"知力"},
	"隐匿":   {"灵巧", "社会", "运势"},
	"强运":   {"运势"},
}

var emokloreAttrParent2 = map[string][]string{
	"治疗": {"知力"},
	"复苏": {"知力", "精神"},
}

var emokloreAttrParent3 = map[string][]string{
	"调查": {"灵巧"},
	"知觉": {"五感"},
	"交涉": {"魅力"},
	"知识": {"知力"},
	"信息": {"社会"},
	"运动": {"身体"},
	"格斗": {"身体"},
	"投掷": {"灵巧"},
	"生存": {"身体"},
	"自我": {"精神"},
	"手工": {"灵巧"},
	"幸运": {"运势"},
}

type singleRoulette struct {
	Name string
	Face int64
	Time int
	Pool []int
}

func pingWebsite(url string) (string, error) {
	chcpCmd := exec.Command("cmd", "/C", "chcp 65001")
	chcpCmd.Run()
	cmd := exec.Command("ping", "-w", "250", url)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

var rouletteMap SyncMap[string, singleRoulette]

func RegisterBuiltinExtFun(self *Dice) {
	aliasHelp := ".alias <别名> <指令> // 将 .&<别名> 定义为指定指令的快捷触发方式\n" +
		".alias --my <别名> <指令> // 将 .&<别名> 定义为个人快捷指令\n" +
		".alias del/rm <别名> // 删除群快捷指令\n" +
		".alias del/rm --my <别名> // 删除个人快捷指令\n" +
		".alias show/list // 显示目前可用的快捷指令\n" +
		".alias help // 查看帮助\n" +
		"// 执行快捷命令见 .& 命令"
	cmdAlias := CmdItemInfo{
		Name:      "alias",
		ShortHelp: aliasHelp,
		Help:      "可以定义一条指令的快捷方式。\n" + aliasHelp,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if len(cmdArgs.Args) == 0 {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			_isPersonal := cmdArgs.GetKwarg("my")
			isPersonal := ctx.MessageType == "private" || _isPersonal != nil

			playerAttrs := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Player.UserID))
			groupAttrs := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Group.GroupID))
			subCmd := cmdArgs.GetArgN(1)

		subParse:
			switch subCmd {
			case "help":
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			case "del", "rm":
				name := cmdArgs.GetArgN(2)
				key := "$g:alias:" + name
				m := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Group.GroupID))
				VarSetValueStr(ctx, "$t指令来源", "群")
				if isPersonal {
					key = "$m:alias:" + name
					m = playerAttrs
					VarSetValueStr(ctx, "$t指令来源", "个人")
				}
				if cmd, ok := m.LoadX(key); ok {
					if cmd != nil && cmd.TypeId == ds.VMTypeString {
						VarSetValueStr(ctx, "$t快捷指令名", name)
						VarSetValueStr(ctx, "$t旧指令", cmd.Value.(string))
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_删除"))
					}
					m.Delete(key)
				} else {
					VarSetValueStr(ctx, "$t快捷指令名", name)
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_删除_未定义"))
				}
			case "list", "show":
				var personCmds, groupCmds []string
				playerAttrs.Range(func(key string, value *ds.VMValue) bool {
					if strings.HasPrefix(key, "$m:alias:") {
						_cmd := key[len("$m:alias:"):]
						if value.TypeId == ds.VMTypeString {
							VarSetValueStr(ctx, "$t快捷指令名", _cmd)
							VarSetValueStr(ctx, "$t指令", value.ToString())
							VarSetValueStr(ctx, "$t指令来源", "个人")
							personCmds = append(personCmds, DiceFormatTmpl(ctx, "核心:快捷指令_列表_单行"))
						}
					}
					return true
				})

				if ctx.MessageType == "group" {
					groupAttrs.Range(func(key string, value *ds.VMValue) bool {
						if strings.HasPrefix(key, "$g:alias:") {
							_cmd := key[len("$g:alias:"):]
							if value.TypeId == ds.VMTypeString {
								VarSetValueStr(ctx, "$t快捷指令名", _cmd)
								VarSetValueStr(ctx, "$t指令", value.ToString())
								VarSetValueStr(ctx, "$t指令来源", "群")
								groupCmds = append(groupCmds, DiceFormatTmpl(ctx, "核心:快捷指令_列表_单行"))
							}
						}

						return false
					})
				}
				sep := DiceFormatTmpl(ctx, "核心:快捷指令_列表_分隔符")
				// 保证群在前个人在后的顺序
				var totalCmds []string
				totalCmds = append(totalCmds, groupCmds...)
				totalCmds = append(totalCmds, personCmds...)
				if len(totalCmds) > 0 {
					VarSetValueStr(ctx, "$t列表内容", strings.Join(totalCmds, sep))
				}

				if len(totalCmds) == 0 {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_列表_空"))
				} else {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_列表"))
				}
			default:
				if len(cmdArgs.Args) < 2 {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_新增_无指令"))
					break
				}
				name := subCmd
				if len(cmdArgs.Args) >= 2 {
					targetCmd := cmdArgs.GetArgN(2)
					for _, prefix := range ctx.Session.Parent.CommandPrefix {
						// 这里依然拦截不了先定义了快捷指令，后添加了新的指令前缀导致出现递归的情况，但是一是这种情况少，二是后面执行阶段也有拦截所以问题不大
						if targetCmd == prefix+"a" || targetCmd == prefix+"&" {
							ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_禁止"))
							break subParse
						}
					}
				}
				_args := cmdArgs.Args[1:]
				for _, kwa := range cmdArgs.Kwargs {
					if kwa.Name != "my" {
						_args = append(_args, kwa.String())
					}
				}
				cmd := strings.TrimSpace(strings.Join(_args, " "))

				m := groupAttrs
				key := "$g:alias:" + name
				VarSetValueStr(ctx, "$t指令来源", "群")
				if isPersonal {
					key = "$m:alias:" + name
					m = playerAttrs
					VarSetValueStr(ctx, "$t指令来源", "个人")
				}

				if oldCmd, ok := m.LoadX(key); ok {
					if oldCmd.TypeId == ds.VMTypeString {
						m.Store(key, ds.NewStrVal(cmd))
						VarSetValueStr(ctx, "$t快捷指令名", name)
						VarSetValueStr(ctx, "$t指令", cmd)
						VarSetValueStr(ctx, "$t旧指令", oldCmd.Value.(string))
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_替换"))
					} else {
						// 防止错误的数据一直卡着
						m.Delete(key)
					}
				} else {
					m.Store(key, ds.NewStrVal(cmd))
					VarSetValueStr(ctx, "$t快捷指令名", name)
					VarSetValueStr(ctx, "$t指令", cmd)
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令_新增"))
				}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	aHelp := ".&/a <快捷指令名> [参数] // 执行对应快捷指令\n" +
		".& help // 查看帮助\n" +
		"// 定义快捷指令见 .alias 命令"
	cmdA := CmdItemInfo{
		Name:      "&",
		ShortHelp: aHelp,
		Help:      "执行一条快捷指令。\n" + aHelp,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if len(cmdArgs.Args) == 0 {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			name := cmdArgs.GetArgN(1)
			if name == "help" {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			log := self.Logger
			args := cmdArgs.Args
			for _, kwa := range cmdArgs.Kwargs {
				args = append(args, kwa.String())
			}

			if msg.MessageType == "group" {
				groupAttrs := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Group.GroupID))
				cmdValue, ok := groupAttrs.LoadX("$g:alias:" + name)
				if ok {
					if cmdValue != nil && cmdValue.TypeId == ds.VMTypeString {
						args[0] = cmdValue.Value.(string)
						targetCmd := strings.Join(args, " ")
						targetArgs := CommandParse(targetCmd, []string{}, self.CommandPrefix, msg.Platform, false)
						if targetArgs != nil {
							log.Infof("群快捷指令映射: .&%s -> %s", cmdArgs.CleanArgs, targetCmd)
							if targetArgs.Command == "a" || targetArgs.Command == "&" {
								return CmdExecuteResult{Matched: true, Solved: true}
							}

							VarSetValueStr(ctx, "$t指令来源", "群")
							VarSetValueStr(ctx, "$t目标指令", targetCmd)
							ctx.AliasPrefixText = DiceFormatTmpl(ctx, "核心:快捷指令触发_前缀")

							ctx.EndPoint.TriggerCommand(ctx, msg, targetArgs)
							return CmdExecuteResult{Matched: true, Solved: true}
						}
					}
				}
			}

			playerAttrs := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Player.UserID))
			cmdValue, ok := playerAttrs.LoadX("$m:alias:" + name)
			if ok {
				if cmdValue != nil && cmdValue.TypeId == ds.VMTypeString {
					args[0] = cmdValue.Value.(string)
					targetCmd := strings.Join(args, " ")
					msg.Message = targetCmd
					targetArgs := CommandParse(targetCmd, []string{}, self.CommandPrefix, msg.Platform, false)
					if targetArgs != nil {
						log.Infof("个人快捷指令映射: .&%s -> %s", cmdArgs.CleanArgs, targetCmd)
						if targetArgs.Command == "a" || targetArgs.Command == "&" {
							return CmdExecuteResult{Matched: true, Solved: true}
						}

						VarSetValueStr(ctx, "$t指令来源", "个人")
						VarSetValueStr(ctx, "$t目标指令", targetCmd)
						ctx.AliasPrefixText = DiceFormatTmpl(ctx, "核心:快捷指令触发_前缀")

						ctx.EndPoint.TriggerCommand(ctx, msg, targetArgs)
						return CmdExecuteResult{Matched: true, Solved: true}
					}
				}
			}

			VarSetValueStr(ctx, "$t目标指令名", name)
			ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:快捷指令触发_无指令"))
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	cmdPing := CmdItemInfo{
		Name:      "ping",
		ShortHelp: ".ping <网站名称> // 触发发送一条回复",
		Help:      "触发回复:\n触发发送一条回复。特别地，如果是qq官方bot，并且是在频道中触发，会以私信消息形式回复",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {

			val := cmdArgs.GetArgN(1)
			switch strings.ToLower(val) {
			case "help":
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			case "baidu":
				ReplyToSender(ctx, msg, "正在向目标网站发起请求")
				pingReturn, _ := pingWebsite("www.baidu.com")
				time.Sleep(2 * time.Second)
				VarSetValueStr(ctx, "$t请求结果", pingReturn)
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "其它:ping响应"))
			case "github":
				ReplyToSender(ctx, msg, "正在向目标网站发起请求")
				pingReturn, _ := pingWebsite("www.github.com")
				time.Sleep(2 * time.Second)
				VarSetValueStr(ctx, "$t请求结果", pingReturn)
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "其它:ping响应"))
			default:
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	cmdSend := CmdItemInfo{
		Name:      "send",
		ShortHelp: ".send // 向骰主留言",
		Help: "留言指令:\n.send XXXXXX // 向骰主留言\n" +
			".send to <对方ID> 要说的话 // 骰主回复，举例. send to QQ:12345 感谢留言\n" +
			".send to <群组ID> 要说的话 // 举例. send to QQ-Group:12345 感谢留言\n" +
			"> 指令.userid可以查看当前群的ID",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			val := cmdArgs.GetArgN(1)
			if val == "to" { //nolint:nestif // TODO
				if ctx.PrivilegeLevel >= 100 {
					uid := cmdArgs.GetArgN(2)
					txt := cmdArgs.GetRestArgsFrom(3)
					if uid != "" && strings.HasPrefix(uid, ctx.EndPoint.Platform) && txt != "" {
						isGroup := strings.Contains(uid, "-Group:")
						txt = fmt.Sprintf("本消息由骰主<%s>通过指令发送:\n", ctx.Player.Name) + txt
						if isGroup {
							ReplyGroup(ctx, &Message{GroupID: uid}, txt)
						} else {
							ReplyPerson(ctx, &Message{Sender: SenderBase{UserID: uid}}, txt)
						}
						ReplyToSender(ctx, msg, "信息已经发送至"+uid)
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}
				ReplyToSender(ctx, msg, "你不具备Master权限")
			} else if val == "help" || val == "" {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			} else {
				if self.MailEnable {
					_ = ctx.Dice.SendMail(cmdArgs.CleanArgs, MailTypeSendNote)
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:留言_已记录"))
					return CmdExecuteResult{Matched: true, Solved: true}
				}
				for _, uid := range ctx.Dice.DiceMasters {
					text := ""

					if ctx.IsCurGroupBotOn {
						text += fmt.Sprintf("一条来自群组<%s>(%s)，作者<%s>(%s)的留言:\n", ctx.Group.GroupName, ctx.Group.GroupID, ctx.Player.Name, ctx.Player.UserID)
					} else {
						text += fmt.Sprintf("一条来自私聊，作者<%s>(%s)的留言:\n", ctx.Player.Name, ctx.Player.UserID)
					}

					text += cmdArgs.CleanArgs
					if strings.Contains(uid, "Group") {
						ctx.EndPoint.Adapter.SendToGroup(ctx, uid, text, "")
					} else {
						ctx.EndPoint.Adapter.SendToPerson(ctx, uid, text, "")
					}
				}
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:留言_已记录"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
		},
	}

	botWelcomeHelp := ".welcome on // 开启\n" +
		".welcome off // 关闭\n" +
		".welcome show // 查看当前欢迎语\n" +
		".welcome set <欢迎语> // 设定欢迎语\n" +
		".welcome clr // 设定欢迎语"
	cmdWelcome := CmdItemInfo{
		Name:              "welcome",
		ShortHelp:         botWelcomeHelp,
		Help:              "新人入群自动发言设定:\n" + botWelcomeHelp,
		DisabledInPrivate: true,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			pRequired := 50 // 50管理 60群主 100master
			if ctx.PrivilegeLevel < pRequired {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			if cmdArgs.IsArgEqual(1, "on") {
				ctx.Group.ShowGroupWelcome = true
				ctx.Group.UpdatedAtTime = time.Now().Unix()
				ReplyToSender(ctx, msg, "入群欢迎语已打开")
			} else if cmdArgs.IsArgEqual(1, "off") {
				ctx.Group.ShowGroupWelcome = false
				ctx.Group.UpdatedAtTime = time.Now().Unix()
				ReplyToSender(ctx, msg, "入群欢迎语已关闭")
			} else if cmdArgs.IsArgEqual(1, "show") {
				welcome := ctx.Group.GroupWelcomeMessage
				var info string
				if ctx.Group.ShowGroupWelcome {
					info = "\n状态: 开启"
				} else {
					info = "\n状态: 关闭"
				}
				ReplyToSender(ctx, msg, "当前欢迎语:\n"+welcome+info)
			} else if cmdArgs.IsArgEqual(1, "clr") {
				ctx.Group.GroupWelcomeMessage = ""
				ctx.Group.ShowGroupWelcome = false
				ctx.Group.UpdatedAtTime = time.Now().Unix()
				ReplyToSender(ctx, msg, "入群欢迎语已清空并关闭")
			} else if _, ok := cmdArgs.EatPrefixWith("set"); ok {
				text2 := strings.TrimSpace(cmdArgs.RawArgs[len("set"):])
				ctx.Group.GroupWelcomeMessage = text2
				ctx.Group.ShowGroupWelcome = true
				ctx.Group.UpdatedAtTime = time.Now().Unix()
				ReplyToSender(ctx, msg, "当前欢迎语设定为:\n"+text2+"\n入群欢迎语已自动打开(注意，会在bot off时起效)")
			} else {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	cmdGugu := CmdItemInfo{
		Name:      "gugu",
		ShortHelp: ".gugu 来源 // 获取一个随机的咕咕理由，带上来源可以看作者",
		Help:      "人工智能鸽子:\n.gugu 来源 // 获取一个随机的咕咕理由，带上来源可以看作者\n.text // 文本指令",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			// p := getPlayerInfoBySender(session, msg)
			isShowFrom := cmdArgs.IsArgEqual(1, "from", "showfrom", "来源", "作者")

			reason := DiceFormatTmpl(ctx, "娱乐:鸽子理由")
			reasonInfo := strings.SplitN(reason, "|", 2)

			text := "🕊️: " + reasonInfo[0]
			if isShowFrom && len(reasonInfo) == 2 {
				text += "\n    ——" + reasonInfo[1]
			}
			ReplyToSender(ctx, msg, text)
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	cmdJrrp := CmdItemInfo{
		Name:      "jrrp",
		ShortHelp: ".jrrp 获得一个D100随机值，一天内不会变化",
		Help:      "今日人品:\n.jrrp 获得一个D100随机值，一天内不会变化",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			rpSeed := (time.Now().Unix() + (8 * 60 * 60)) / (24 * 60 * 60)
			rpSeed += int64(fingerprint(ctx.EndPoint.UserID))
			rpSeed += int64(fingerprint(ctx.Player.UserID))
			src := rand.NewPCG(uint64(rpSeed), uint64(-rpSeed))
			randItem := rand.New(src)
			rp := randItem.Int64()%100 + 1

			VarSetValueInt64(ctx, "$t人品", rp)
			ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "娱乐:今日人品"))
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	helpDuel := ".duel //和骰子决斗"
	cmdDuel := CmdItemInfo{
		Name:      "duel",
		ShortHelp: helpDuel,
		Help:      "决斗:\n" + helpDuel,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			//duelPrefix := DiceFormatTmpl(ctx,"决斗前置文本")
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	textHelp := ".text <文本模板> // 文本指令，例: .text 看看手气: {1d16}"
	cmdText := CmdItemInfo{
		Name:      "text",
		ShortHelp: textHelp,
		Help:      "文本模板指令:\n" + textHelp,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if ctx.Dice.TextCmdTrustOnly {
				// 检查master和信任权限
				// 拒绝无权限访问
				if ctx.PrivilegeLevel < 70 {
					ReplyToSender(ctx, msg, "你不具备Master权限")
					return CmdExecuteResult{Matched: true, Solved: true}
				}
			}
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			tmpl := ctx.Group.GetCharTemplate(ctx.Dice)
			ctx.Eval(tmpl.PreloadCode, nil)
			val := cmdArgs.GetArgN(1)

			if val != "" {
				ctx.Player.TempValueAlias = nil // 防止dnd的hp被转为“生命值”
				r, _, err := DiceExprTextBase(ctx, cmdArgs.CleanArgs, RollExtraFlags{DisableBlock: false, V2Only: true})

				if err == nil {
					text := r.ToString()

					if kw := cmdArgs.GetKwarg("asm"); r != nil && kw != nil {
						if ctx.PrivilegeLevel >= 40 {
							asm := r.GetAsmText()
							text += "\n" + asm
						}
					}

					if r.legacy != nil {
						text += "\n" + "* 当前表达式在RollVM V2中无法报错，建议修改：" + r.vm.Error.Error()
					}

					seemsCommand := false
					if strings.HasPrefix(text, ".") || strings.HasPrefix(text, "。") || strings.HasPrefix(text, "!") || strings.HasPrefix(text, "/") {
						seemsCommand = true
						if strings.HasPrefix(text, "..") || strings.HasPrefix(text, "。。") || strings.HasPrefix(text, "!!") {
							seemsCommand = false
						}
					}

					if seemsCommand {
						ReplyToSender(ctx, msg, "你可能在利用text让骰子发出指令文本，这被视为恶意行为并已经记录")
					} else {
						ReplyToSender(ctx, msg, text)
					}
				} else {
					ReplyToSender(ctx, msg, "执行出错:"+err.Error())
				}
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
		},
	}

	self.RegisterExtension(&ExtInfo{
		Name:            "fun", // 扩展的名称，需要用于指令中，写简短点      2024.05.10: 目前被看成是 function 的缩写了（
		Version:         "1.1.0",
		Brief:           "功能扩展，主要提供快捷指令、ping、welcome等额外指令，同时也包括今日人品、智能鸽子等娱乐相关指令。同时，小众规则指令暂时也放在本扩展中",
		AutoActive:      true, // 是否自动开启
		ActiveOnPrivate: true,
		Author:          "木落",
		Official:        true,
		OnCommandReceived: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) {
		},
		OnLoad: func() {
		},
		GetDescText: GetExtensionDesc,
		CmdMap: CmdMapCls{
			"alias":   &cmdAlias,
			"&":       &cmdA,
			"a":       &cmdA,
			"ping":    &cmdPing,
			"send":    &cmdSend,
			"welcome": &cmdWelcome,
			"gugu":    &cmdGugu,
			"咕咕":      &cmdGugu,
			"jrrp":    &cmdJrrp,
			"duel":    &cmdDuel,
			"text":    &cmdText,
		},
	})
}

func fingerprint(b string) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(b))
	return hash.Sum64()
}
