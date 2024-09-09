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
		Help:      "黑名单接收:\n" + HelpForShikiBot,
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
					if !(msg.Platform == "QQ-CH" || ctx.Dice.BotExtFreeSwitch || ctx.PrivilegeLevel >= 40) {
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
		},
	})
}