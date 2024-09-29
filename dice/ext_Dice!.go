package dice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type warningMessage struct {
	Wid       int64  `json:"wid"`
	Type      string `json:"type"`
	Danger    int    `json:"danger"`
	FromGroup int64  `json:"fromGroup"`
	FromGID   int64  `json:"fromGID"`
	FromQQ    int64  `json:"fromQQ"`
	FromUID   int64  `json:"fromUID"`
	InviterQQ int64  `json:"inviterQQ"`
	Time      string `json:"time"`
	Note      string `json:"note"`
	DiceMaid  int64  `json:"DiceMaid"`
	MasterQQ  int64  `json:"masterQQ"`
	Comment   string `json:"comment"`
}

func RegisterBuiltinShikiCommands(d *Dice) {
	helpForShikiWarning := "溯洄warning播报处理"
	cmdShikiWarning := CmdItemInfo{
		Name:      "Dice!warning",
		ShortHelp: helpForShikiWarning,
		Help:      "黑名单接收:\n" + helpForShikiWarning,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if ctx.PrivilegeLevel < 100 {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			if strings.ToLower(cmdArgs.GetArgN(1)) == "generate" || strings.ToLower(cmdArgs.GetArgN(1)) == "gene" {
				subval := strings.ToLower(cmdArgs.GetArgN(2))
				switch subval {
				case "qq":
					BlackType := strings.ToLower(cmdArgs.GetArgN(3))
					BlackQQ := cmdArgs.GetArgN(4)
					var warningStruct warningMessage
					var warningNote string
					var warningDanger int
					if BlackType == "ban" {
						warningStruct.Type = "ban"
						warningNote = "被禁言"
						warningDanger = 2
					} else if BlackType == "mute" {
						warningStruct.Type = "mute"
						warningNote = "被禁言"
						warningDanger = 2
					} else if BlackType == "kick" {
						warningStruct.Type = "kick"
						warningNote = "被踢出"
						warningDanger = 2
					} else if BlackType == "spam" {
						warningStruct.Type = "spam"
						warningNote = "被标记为刷屏"
						warningDanger = 1
					} else {
						warningStruct.Type = "other"
						warningNote = "其他原因"
						warningDanger = 1
					}
					warningStruct.Time = time.Now().Format("2006-01-02 15:04:05")
					tmpvar, _ := strconv.Atoi(ctx.EndPoint.UserID)
					warningStruct.DiceMaid = int64(tmpvar)
					warningStruct.Comment = fmt.Sprintf("%s%s%s%s%s%s", warningStruct.Time, "由骰主: ", ctx.Player.Name, "于群: ", ctx.Group.GroupName, "中生成")
					tmpvar, _ = strconv.Atoi(ctx.Player.UserID)
					warningStruct.MasterQQ = int64(tmpvar)
					warningStruct.Danger = warningDanger
					warningStruct.Note = warningNote

					tmpvar, _ = strconv.Atoi(BlackQQ)
					warningStruct.FromQQ = int64(tmpvar)
					warningStruct.FromUID = int64(tmpvar)

					warningJson, _ := json.Marshal(warningStruct)
					reply := fmt.Sprintf("%s%s", "!warning", string(warningJson))
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				case "group":
					BlackType := strings.ToLower(cmdArgs.GetArgN(3))
					BlackGroup := cmdArgs.GetArgN(4)
					var warningStruct warningMessage
					var warningNote string
					var warningDanger int
					if BlackType == "ban" {
						warningStruct.Type = "ban"
						warningNote = "被禁言"
						warningDanger = 2
					} else if BlackType == "mute" {
						warningStruct.Type = "mute"
						warningNote = "被禁言"
						warningDanger = 2
					} else if BlackType == "kick" {
						warningStruct.Type = "kick"
						warningNote = "被踢出"
						warningDanger = 2
					} else if BlackType == "spam" {
						warningStruct.Type = "spam"
						warningNote = "被标记为刷屏"
						warningDanger = 1
					} else {
						warningStruct.Type = "other"
						warningNote = "其他原因"
						warningDanger = 1
					}
					warningStruct.Time = time.Now().Format("2006-01-02 15:04:05")
					tmpvar, _ := strconv.Atoi(ctx.EndPoint.UserID)
					warningStruct.DiceMaid = int64(tmpvar)
					warningStruct.Comment = fmt.Sprintf("%s%s%s%s%s%s", warningStruct.Time, "由骰主: ", ctx.Player.Name, "于群: ", ctx.Group.GroupName, "中生成")
					tmpvar, _ = strconv.Atoi(ctx.Player.UserID)
					warningStruct.MasterQQ = int64(tmpvar)
					warningStruct.Danger = warningDanger
					warningStruct.Note = warningNote

					tmpvar, _ = strconv.Atoi(BlackGroup)
					warningStruct.FromGroup = int64(tmpvar)
					warningStruct.FromGID = int64(tmpvar)

					warningJson, _ := json.Marshal(warningStruct)
					reply := fmt.Sprintf("%s%s", "!warning", string(warningJson))
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				default:
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}

			} else {
				// 解析警告信息
				re := regexp.MustCompile(`\s`)
				warningInformation := re.ReplaceAllString(cmdArgs.RawArgs, "")
				var warningStruct warningMessage
				err := json.Unmarshal([]byte(warningInformation), &warningStruct)
				if err != nil {
					ReplyToSender(ctx, msg, "警告信息解析失败"+cmdArgs.RawArgs)
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				retMes := ""
				if warningStruct.Type != "erase" {
					// 处理fromGroup和fromQQ
					if warningStruct.FromGroup != 0 {
						warningEventGroup := fmt.Sprintf("QQ-Group:%d", warningStruct.FromGroup)
						item, ok := d.BanList.GetByID(warningEventGroup)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(warningEventGroup, d.BanList.ThresholdBan, warningStruct.Comment, "溯洄广播黑名单同步", ctx)
							retMes += fmt.Sprintf("已将%s加入黑名单✓\n", warningEventGroup)
						} else {
							retMes += fmt.Sprintf("%s已在黑名单中✓\n", warningEventGroup)
						}
					}
					if warningStruct.FromGID != 0 {
						warningEventGroup := fmt.Sprintf("QQ-Group:%d", warningStruct.FromGID)
						item, ok := d.BanList.GetByID(warningEventGroup)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(warningEventGroup, d.BanList.ThresholdBan, warningStruct.Comment, "溯洄广播黑名单同步", ctx)
							retMes += fmt.Sprintf("已将%s加入黑名单✓\n", warningEventGroup)
						} else {
							retMes += fmt.Sprintf("%s已在黑名单中✓\n", warningEventGroup)
						}

					}
					if warningStruct.FromQQ != 0 {
						warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.FromQQ)
						item, ok := d.BanList.GetByID(warningEventQQ)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(warningEventQQ, d.BanList.ThresholdBan, warningStruct.Comment, "溯洄广播黑名单同步", ctx)
							retMes += fmt.Sprintf("已将%s加入黑名单✓\n", warningEventQQ)
						} else {
							retMes += fmt.Sprintf("%s已在黑名单中✓\n", warningEventQQ)
						}
					}

					if warningStruct.FromUID != 0 {
						warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.FromUID)
						item, ok := d.BanList.GetByID(warningEventQQ)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(warningEventQQ, d.BanList.ThresholdBan, warningStruct.Comment, "溯洄广播黑名单同步", ctx)
							retMes += fmt.Sprintf("已将%s加入黑名单✓\n", warningEventQQ)
						} else {
							retMes += fmt.Sprintf("%s已在黑名单中✓\n", warningEventQQ)
						}
					}
					if warningStruct.InviterQQ != 0 {
						if warningStruct.Type == "ban" {
							warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.InviterQQ)
							item, ok := d.BanList.GetByID(warningEventQQ)
							if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
								d.BanList.AddScoreBase(warningEventQQ, d.BanList.ThresholdBan, warningStruct.Comment, "溯洄广播黑名单同步，邀请人连带", ctx)
								retMes += fmt.Sprintf("已将%s加入黑名单✓\n", warningEventQQ)
							} else {
								retMes += fmt.Sprintf("%s已在黑名单中✓\n", warningEventQQ)
							}
						}
					}
				} else {
					// 处理fromGroup和fromQQ
					if warningStruct.FromGroup != 0 {
						warningEventGroup := fmt.Sprintf("QQ-Group:%d", warningStruct.FromGroup)
						item, ok := d.BanList.GetByID(warningEventGroup)
						if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
							item.Score = 0
							item.Rank = BanRankNormal
							retMes += fmt.Sprintf("已将%s移除黑名单✓\n", warningEventGroup)
						} else {
							retMes += fmt.Sprintf("%s并未在黑名单中✓\n", warningEventGroup)
						}
					}
					if warningStruct.FromGID != 0 {
						warningEventGroup := fmt.Sprintf("QQ-Group:%d", warningStruct.FromGID)
						item, ok := d.BanList.GetByID(warningEventGroup)
						if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
							item.Score = 0
							item.Rank = BanRankNormal
							retMes += fmt.Sprintf("已将%s移除黑名单✓\n", warningEventGroup)
						} else {
							retMes += fmt.Sprintf("%s并未在黑名单中✓\n", warningEventGroup)
						}
					}
					if warningStruct.FromQQ != 0 {
						warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.FromQQ)
						item, ok := d.BanList.GetByID(warningEventQQ)
						if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
							item.Score = 0
							item.Rank = BanRankNormal
							retMes += fmt.Sprintf("已将%s移除黑名单✓\n", warningEventQQ)
						} else {
							retMes += fmt.Sprintf("%s并未在黑名单中✓\n", warningEventQQ)
						}
					}
					if warningStruct.FromUID != 0 {
						warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.FromUID)
						item, ok := d.BanList.GetByID(warningEventQQ)
						if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
							item.Score = 0
							item.Rank = BanRankNormal
							retMes += fmt.Sprintf("已将%s移除黑名单✓\n", warningEventQQ)
						} else {
							retMes += fmt.Sprintf("%s并未在黑名单中✓\n", warningEventQQ)
						}
					}
					if warningStruct.InviterQQ != 0 {
						warningEventQQ := fmt.Sprintf("QQ:%d", warningStruct.InviterQQ)
						item, ok := d.BanList.GetByID(warningEventQQ)
						if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
							item.Score = 0
							item.Rank = BanRankNormal
							retMes += fmt.Sprintf("已将%s移除黑名单✓\n", warningEventQQ)
						} else {
							retMes += fmt.Sprintf("%s并未在黑名单中✓\n", warningEventQQ)
						}
					}

				}
				var warningInformationJson bytes.Buffer
				_ = json.Indent(&warningInformationJson, []byte(warningInformation), "", "    ")

				ReplyToSender(ctx, msg, retMes)
				ReplyToSender(ctx, msg, fmt.Sprintf("%s %s已通知%s不良记录%d:\n!warning%s", time.Now().Format("2006-01-02 15:04:05"), ctx.Player.Name, ctx.EndPoint.Nickname, warningStruct.Wid, warningInformationJson.String()))
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}
	HelpForShikiBot := ".bot on/off"
	cmdShikiBot := CmdItemInfo{
		Name:      "Dice!Bot",
		ShortHelp: HelpForShikiBot,
		Help:      "骰子开关:\n" + HelpForShikiBot,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			inGroup := msg.MessageType == "group"

			if inGroup {
				// 不响应裸指令选项
				if len(cmdArgs.At) < 1 && ctx.Dice.IgnoreUnaddressedBotCmd {
					return CmdExecuteResult{Matched: true, Solved: false}
				}
				// 不响应at其他人
				if cmdArgs.SomeoneBeMentionedButNotMe {
					return CmdExecuteResult{Matched: true, Solved: false}
				}
			}

			if len(cmdArgs.Args) > 0 {
				if cmdArgs.SomeoneBeMentionedButNotMe {
					return CmdExecuteResult{Matched: true, Solved: false}
				}

				cmdArgs.ChopPrefixToArgsWith("on", "off")

				matchNumber := func() (bool, bool) {
					txt := cmdArgs.GetArgN(2)
					if len(txt) >= 4 {
						if strings.HasSuffix(ctx.EndPoint.UserID, txt) {
							return true, txt != ""
						}
					}
					return false, txt != ""
				}

				isMe, exists := matchNumber()
				if exists && !isMe {
					return CmdExecuteResult{Matched: true, Solved: false}
				}
				val := cmdArgs.GetArgN(1)
				switch strings.ToLower(val) {
				case "on":
					if !(ctx.Dice.BotExtFreeSwitch || ctx.PrivilegeLevel >= 40) {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理/邀请者"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					if ctx.IsPrivate {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_私聊不可用"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					SetBotOnAtGroup(ctx, msg.GroupID)
					// TODO：ServiceAtNew此处忽略是否合理？
					ctx.Group, _ = ctx.Session.ServiceAtNew.Load(msg.GroupID)
					ctx.IsCurGroupBotOn = true

					text := DiceFormatTmpl(ctx, "核心:骰子开启")
					ReplyToSender(ctx, msg, text)

					return CmdExecuteResult{Matched: true, Solved: true}
				case "off":
					if !(msg.Platform == "QQ-CH" || ctx.Dice.BotExtFreeSwitch || ctx.PrivilegeLevel >= 40) {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理/邀请者"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					if ctx.IsPrivate {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_私聊不可用"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					SetBotOffAtGroup(ctx, ctx.Group.GroupID)
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰子关闭"))
					return CmdExecuteResult{Matched: true, Solved: true}
				default:
					return CmdExecuteResult{Matched: true, Solved: true}
				}
			}
			if cmdArgs.SomeoneBeMentionedButNotMe {
				return CmdExecuteResult{Matched: false, Solved: false}
			}
			ShikiBotHeader := DiceFormatTmpl(ctx, "核心:bot信息头")
			ver := VERSION.String()
			ShikiBotAuthor := "星界 & 星界之仆"
			text := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s", ShikiBotHeader, " Tempest Dice by ", ShikiBotAuthor, " Ver ", ver, " on Tempest Dice 驱动器", ctx.EndPoint.Platform, " by ", ShikiBotAuthor, " ver ", ver, " for ", ctx.EndPoint.Platform)
			ReplyToSender(ctx, msg, text)

			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}
	HelpForShikiAdminNotceAdd := ""
	HelpForShikiAdminNotceDel := ""
	HelpForShikiAdminNoticeClr := ""
	HelpForShikiAdminNotceList := ""
	HelpForShikiAdmin := "" + HelpForShikiAdminNotceAdd + "\n" + HelpForShikiAdminNotceDel + "\n" + HelpForShikiAdminNoticeClr + "\n" + HelpForShikiAdminNotceList
	HelpForShikiAdminNotce := ""
	cmdShikiAdmin := CmdItemInfo{
		Name:      "Dice!Admin",
		ShortHelp: HelpForShikiAdmin,
		Help:      "骰子管理:\n" + HelpForShikiAdmin,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if ctx.PrivilegeLevel < 100 {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			cmdNum := len(cmdArgs.Args)
			val := cmdArgs.GetArgN(1)
			subval := cmdArgs.GetArgN(2)
			trdval := cmdArgs.GetArgN(3)
			var uid string
			getID := func() string {
				if cmdArgs.IsArgEqual(2, "-") || cmdArgs.IsArgEqual(2, "+") {
					id := cmdArgs.GetArgN(3)
					if id == "" {
						return ""
					}

					isGroup := cmdArgs.IsArgEqual(1, "blackgroup")
					return FormatDiceID(ctx, id, isGroup)
				}

				arg := cmdArgs.GetArgN(2)
				if !strings.Contains(arg, ":") {
					return ""
				}
				return arg
			}
			switch strings.ToLower(val) {
			case "notice":
				switch strings.ToLower(subval) {
				case "help":
					ReplyToSender(ctx, msg, HelpForShikiAdminNotce)
					return CmdExecuteResult{Matched: true, Solved: true}
				case "+", "add":
					if cmdNum < 4 {
						if trdval != "" {
							if strings.ToLower(trdval) == "help" {
								ReplyToSender(ctx, msg, HelpForShikiAdminNotceAdd)
								return CmdExecuteResult{Matched: true, Solved: true}
							}
							if strings.HasPrefix(trdval, "g") {
								trdval = strings.ReplaceAll(trdval, "g", "QQ-Group:")
							} else if strings.HasPrefix(trdval, "p") {
								trdval = strings.ReplaceAll(trdval, "p", "QQ:")
							}
							if strings.HasPrefix(trdval, "QQ:") || strings.HasPrefix(trdval, "QQ-Group:") {
								d.NoticeIDs = append(d.NoticeIDs, trdval)
								d.Save(false)
								ReplyToSender(ctx, msg, "骰子管理: 已添加通知窗口"+trdval)
							}
						}
					} else {
						ReplyToSender(ctx, msg, HelpForShikiAdminNotceAdd)
						return CmdExecuteResult{Matched: true, Solved: true}
					}
				case "-", "del":
					if cmdNum >= 4 {
						if strings.HasPrefix(trdval, "g") {
							trdval = strings.ReplaceAll(trdval, "g", "QQ-Group:")
						} else if strings.HasPrefix(trdval, "p") {
							trdval = strings.ReplaceAll(trdval, "p", "QQ:")
						}
						if strings.HasPrefix(trdval, "QQ:") || strings.HasPrefix(trdval, "QQ-Group:") {
							if len(d.NoticeIDs) == 0 {
								ReplyToSender(ctx, msg, "骰子管理: 通知列表为空")
								return CmdExecuteResult{Matched: true, Solved: true}
							} else if trdval == "help" {
								ReplyToSender(ctx, msg, HelpForShikiAdminNotceDel)
								return CmdExecuteResult{Matched: true, Solved: true}
							}
							for i, id := range d.NoticeIDs {
								if id == trdval {
									d.NoticeIDs = append(d.NoticeIDs[:i], d.NoticeIDs[i+1:]...)
									d.Save(false)
									ReplyToSender(ctx, msg, "骰子管理: 已删除通知窗口"+trdval)
									break
								}
							}
						}
					} else {
						ReplyToSender(ctx, msg, HelpForShikiAdminNotceDel)
						return CmdExecuteResult{Matched: true, Solved: true}
					}

				case "clr":
					if strings.ToLower(trdval) == "help" {
						ReplyToSender(ctx, msg, HelpForShikiAdminNoticeClr)
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					d.NoticeIDs = []string{}
					d.Save(false)
					ReplyToSender(ctx, msg, "骰子管理: 通知列表已清空")
					return CmdExecuteResult{Matched: true, Solved: true}
				case "list":
					if len(d.NoticeIDs) == 0 {
						ReplyToSender(ctx, msg, "骰子管理: 通知列表为空")
						return CmdExecuteResult{Matched: true, Solved: true}
					} else if strings.ToLower(trdval) == "help" {
						ReplyToSender(ctx, msg, "")
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					text := "骰子管理: 通知列表\n"
					for _, id := range d.NoticeIDs {
						text += id + "\n"
					}
					ReplyToSender(ctx, msg, text)
					return CmdExecuteResult{Matched: true, Solved: true}
				default:
					ReplyToSender(ctx, msg, HelpForShikiAdminNotceList)
					return CmdExecuteResult{Matched: true, Solved: true}
				}
			case "dismiss":
				gid := cmdArgs.GetArgN(2)
				if gid == "" || gid == "help" {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}

				n := strings.Split(gid, ":") // 不验证是否合法，反正下面会检查是否在 ServiceAtNew
				if strings.HasPrefix(gid, "g") {
					gid = strings.ReplaceAll(gid, "g", "")
				}
				gid = "QQ-Group:" + gid // 强制当作QQ群聊处理
				gp, ok := ctx.Session.ServiceAtNew.Load(gid)
				if !ok || len(n[0]) < 2 {
					ReplyToSender(ctx, msg, fmt.Sprintf("群组列表中没有找到%s", gid))
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				// 既然是骰主自己操作，就不通知了
				// 除非有多骰主……
				ReplyToSender(ctx, msg, fmt.Sprintf("收到指令，将在5秒后退出群组%s", gp.GroupID))

				txt := "注意，收到骰主指令，5秒后将从该群组退出。"
				wherefore := cmdArgs.GetArgN(3)
				if wherefore != "" {
					txt += fmt.Sprintf("原因: %s", wherefore)
				}

				ReplyGroup(ctx, &Message{GroupID: gp.GroupID}, txt)

				mctx := &MsgContext{
					MessageType: "group",
					Group:       gp,
					EndPoint:    ctx.EndPoint,
					Session:     ctx.Session,
					Dice:        ctx.Dice,
					IsPrivate:   false,
				}
				// SetBotOffAtGroup(mctx, gp.GroupID)
				time.Sleep(3 * time.Second)
				gp.DiceIDExistsMap.Delete(mctx.EndPoint.UserID)
				gp.UpdatedAtTime = time.Now().Unix()
				mctx.EndPoint.Adapter.QuitGroup(mctx, gp.GroupID)

				return CmdExecuteResult{Matched: true, Solved: true}

			case "ban", "blk", "black":
				if len(cmdArgs.Args) < 3 {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				} else if trdval == "help" {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				} else {
					if strings.HasPrefix(subval, "g") {
						subval = strings.ReplaceAll(subval, "g", "QQ-Group:")
					} else if strings.HasPrefix(subval, "p") {
						subval = strings.ReplaceAll(subval, "p", "QQ:")
					}
					if strings.HasPrefix(subval, "QQ:") || strings.HasPrefix(subval, "QQ-Group:") {
						uid := subval
						var reason string
						if trdval == "" {
							reason = "骰主指令"
						} else {
							reason = trdval
						}
						d.BanList.AddScoreBase(uid, d.BanList.ThresholdBan, "骰主指令", reason, ctx)
						ReplyToSender(ctx, msg, fmt.Sprintf("已将用户 %s 加入黑名单，原因: %s", uid, reason))
					}
				}
			case "blackqq":
				var subval = cmdArgs.GetArgN(2)
				if subval == "-" {
					uid = getID()
					if uid == "" {
						return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
					}

					item, ok := d.BanList.GetByID(uid)
					if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
						ReplyToSender(ctx, msg, "找不到用户")
						break
					}

					ReplyToSender(ctx, msg, fmt.Sprintf("已将用户 %s 移出 %s 列表", uid, BanRankText[item.Rank]))
					item.Score = 0
					item.Rank = BanRankNormal

				} else if subval == "+" {
					uid = getID()
					if uid == "" {
						return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
					}
					reason := cmdArgs.GetArgN(4)
					if reason == "" {
						reason = "骰主指令"
					}
					d.BanList.AddScoreBase(uid, d.BanList.ThresholdBan, "骰主指令", reason, ctx)
					ReplyToSender(ctx, msg, fmt.Sprintf("已将用户 %s 加入黑名单，原因: %s", uid, reason))

				} else {
					return CmdExecuteResult{Matched: true, Solved: false, ShowHelp: true}
				}
			case "blackgroup":
				var subval = cmdArgs.GetArgN(2)
				if subval == "-" {
					uid = getID()
					if uid == "" {
						return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
					}

					item, ok := d.BanList.GetByID(uid)
					if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
						ReplyToSender(ctx, msg, "找不到群组")
						break
					}

					ReplyToSender(ctx, msg, fmt.Sprintf("已将群组 %s 移出%s列表", uid, BanRankText[item.Rank]))
					item.Score = 0
					item.Rank = BanRankNormal

				} else if subval == "+" {
					uid = getID()
					if uid == "" {
						return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
					}
					reason := cmdArgs.GetArgN(4)
					if reason == "" {
						reason = "骰主指令"
					}
					d.BanList.AddScoreBase(uid, d.BanList.ThresholdBan, "骰主指令", reason, ctx)
					ReplyToSender(ctx, msg, fmt.Sprintf("已将群组 %s 加入黑名单，原因: %s", uid, reason))

				} else {
					return CmdExecuteResult{Matched: true, Solved: false, ShowHelp: true}
				}

			default:
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}
	d.RegisterExtension(&ExtInfo{
		Name:            "Dice!Core", // 扩展的名称，需要用于指令中，写简短点      2024.05.10: 目前被看成是 function 的缩写了（
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
			"warning": &cmdShikiWarning,
			"bot":     &cmdShikiBot,
			"admin":   &cmdShikiAdmin,
		},
	})
}
