package dice

import (
	"fmt"
	"strings"

	"github.com/juliangruber/go-intersect"
)

/** 这一条指令不能移除 */
func (d *Dice) registerExtLoader() {
	helpExt := ".ext // 查看扩展列表"
	cmdExt := &CmdItemInfo{
		Name:      "ext",
		ShortHelp: helpExt,
		Help:      "群扩展模块管理:\n" + helpExt,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			showList := func() {
				text := "检测到以下扩展(名称-版本-作者)：\n"
				for index, i := range ctx.Dice.ExtList {
					state := "关"
					for _, j := range ctx.Group.ActivatedExtList {
						if i.Name == j.Name {
							state = "开"
							break
						}
					}
					var officialMark string
					if i.Official {
						officialMark = "[官方]"
					}
					author := i.Author
					if author == "" {
						author = "<未注明>"
					}
					aliases := ""
					if len(i.Aliases) > 0 {
						aliases = "(" + strings.Join(i.Aliases, ",") + ")"
					}
					text += fmt.Sprintf("%d. [%s]%s%s %s - %s - %s\n", index+1, state, officialMark, i.Name, aliases, i.Version, author)
				}
				text += "使用命令: .ext <扩展名> on/off 可以在当前群开启或关闭某扩展。\n"
				text += "命令: .ext <扩展名> 可以查看扩展介绍及帮助"
				ReplyToSender(ctx, msg, text)
			}

			if len(cmdArgs.Args) == 0 {
				showList()
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			var last int
			if len(cmdArgs.Args) >= 2 {
				last = len(cmdArgs.Args)
			}

			//nolint:nestif
			if cmdArgs.IsArgEqual(1, "list") {
				showList()
			} else if cmdArgs.IsArgEqual(last, "on") {
				if !ctx.Dice.BotExtFreeSwitch && ctx.PrivilegeLevel < 40 {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理/邀请者"))
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				checkConflict := func(ext *ExtInfo) []string {
					var actived []string
					for _, i := range ctx.Group.ActivatedExtList {
						actived = append(actived, i.Name)
					}

					if ext.ConflictWith != nil {
						var ret []string
						for _, i := range intersect.Simple(actived, ext.ConflictWith) {
							ret = append(ret, i.(string))
						}
						return ret
					}
					return []string{}
				}

				var extNames []string
				var conflictsAll []string
				for index := 0; index < len(cmdArgs.Args); index++ {
					extName := strings.ToLower(cmdArgs.Args[index])
					if i := d.ExtFind(extName); i != nil {
						extNames = append(extNames, extName)
						conflictsAll = append(conflictsAll, checkConflict(i)...)
						ctx.Group.ExtActive(i)
					}
				}

				if len(extNames) == 0 {
					ReplyToSender(ctx, msg, "输入的扩展类别名无效")
				} else {
					text := fmt.Sprintf("打开扩展 %s", strings.Join(extNames, ","))
					if len(conflictsAll) > 0 {
						text += "\n检测到可能冲突的扩展，建议关闭: " + strings.Join(conflictsAll, ",")
						text += "\n对于扩展中存在的同名指令，则越晚开启的扩展，优先级越高。"
					}
					ReplyToSender(ctx, msg, text)
				}
			} else if cmdArgs.IsArgEqual(last, "off") {
				if !ctx.Dice.BotExtFreeSwitch && ctx.PrivilegeLevel < 40 {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理/邀请者"))
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				var closed []string
				var notfound []string
				for index := 0; index < len(cmdArgs.Args); index++ {
					extName := cmdArgs.Args[index]
					extName = d.ExtAliasToName(extName)
					ei := ctx.Group.ExtInactiveByName(extName)
					if ei != nil {
						closed = append(closed, ei.Name)
					} else {
						notfound = append(notfound, extName)
					}
				}

				var text string

				if len(closed) > 0 {
					text += fmt.Sprintf("关闭扩展: %s", strings.Join(closed, ","))
				} else {
					text += fmt.Sprintf(" 已关闭或未找到: %s", strings.Join(notfound, ","))
				}
				ReplyToSender(ctx, msg, text)
				return CmdExecuteResult{Matched: true, Solved: true}
			} else {
				extName := cmdArgs.Args[0]
				if i := d.ExtFind(extName); i != nil {
					text := fmt.Sprintf("> [%s] 版本%s 作者%s\n", i.Name, i.Version, i.Author)
					i.callWithJsCheck(d, func() {
						ReplyToSender(ctx, msg, text+i.GetDescText(i))
					})
					return CmdExecuteResult{Matched: true, Solved: true}
				}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}
	d.CmdMap["ext"] = cmdExt
}
