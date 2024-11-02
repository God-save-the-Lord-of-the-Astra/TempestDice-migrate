package dice

import "strings"

func RegisterBuiltinGithubSubscribeCommands(d *Dice) {
	helpForGithubSubscribeAdd := ".github add <名称> <订阅地址> 添加一个订阅"
	helpForGithubSubscribeRm := ".github rm <名称> 删除一个订阅"
	helpForGithubSubscribeBrief := "Github订阅模块"
	helpForGithubSubscribe := helpForGithubSubscribeBrief + "\n" + helpForGithubSubscribeAdd
	cmdShikGithubSubscribe := CmdItemInfo{
		Name:      "Dice!warning",
		ShortHelp: helpForGithubSubscribe,
		Help:      "黑名单接收:\n" + helpForGithubSubscribe,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if ctx.PrivilegeLevel < 100 {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			if len(cmdArgs.Args) < 2 {
				ReplyToSender(ctx, msg, helpForGithubSubscribe)
				return CmdExecuteResult{Matched: true, Solved: true}
			} else {
				val := cmdArgs.GetArgN(1)
				subval := cmdArgs.GetArgN(2)
				switch strings.ToLower(val) {
				case "add", "+":
					if len(cmdArgs.Args) != 3 {
						ReplyToSender(ctx, msg, helpForGithubSubscribeAdd)
						return CmdExecuteResult{Matched: true, Solved: true}
					} else {
						trdval := cmdArgs.GetArgN(3)
						// 获取订阅信息
						newsub := GitHubSubscribeInfo{
							repoAlias: subval,
							repoUrl:   trdval,
							lastSHA:   "",
						}
						//newsub.Update()
						d.GitHubSubscribeList = append(d.GitHubSubscribeList, newsub)
						d.Save(false)
					}
				case "rm", "-":
					if len(cmdArgs.Args) != 2 {
						ReplyToSender(ctx, msg, helpForGithubSubscribeRm)
						return CmdExecuteResult{Matched: true, Solved: true}
					} else {
						trdval := cmdArgs.GetArgN(2)
						// 获取订阅信息
						for i, v := range d.GitHubSubscribeList {
							if v.repoAlias == trdval {
								d.GitHubSubscribeList = append(d.GitHubSubscribeList[:i], d.GitHubSubscribeList[i+1:]...)
								break
							}
						}
						d.Save(false)
					}
				case "list":
					if len(cmdArgs.Args) != 1 {
						ReplyToSender(ctx, msg, helpForGithubSubscribe)
						return CmdExecuteResult{Matched: true, Solved: true}
					} else {
						var reply string
						if len(d.GitHubSubscribeList) == 0 {
							reply = "无订阅"
						}
						for _, v := range d.GitHubSubscribeList {
							reply += v.repoAlias + " " + v.repoUrl + "\n"
						}
						ReplyToSender(ctx, msg, reply)
					}
				default:
					ReplyToSender(ctx, msg, helpForGithubSubscribe)
					return CmdExecuteResult{Matched: true, Solved: true}
				}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		}}
	d.RegisterExtension(&ExtInfo{
		Name:            "GithubSubscribe", // 扩展的名称，需要用于指令中，写简短点      2024.05.10: 目前被看成是 function 的缩写了（
		Version:         "1.0.0",
		Brief:           "Dice!核心指令",
		AutoActive:      false, // 是否自动开启
		ActiveOnPrivate: true,
		Author:          "海棠、星界之主",
		Official:        true,
		OnCommandReceived: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) {
		},
		OnLoad: func() {
		},
		GetDescText: GetExtensionDesc,
		CmdMap: CmdMapCls{
			"github": &cmdShikGithubSubscribe,
		},
	})
}
