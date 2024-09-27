package dice

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon"
	ds "github.com/sealdice/dicescript"
)

var guguText = `
{$t玩家}为了拯救公主前往了巨龙的巢穴，还没赶回来！|鹊鹊结合实际经历创作
{$t玩家}在来开团的路上被巨龙叼走了！|鹊鹊结合实际经历创作
来的路上出现了哥布林劫匪！{$t玩家}大概是赶不过来了！|鹊鹊结合实际经历创作
咕咕咕~广场上的鸽子把{$t玩家}叼回了巢穴~|鹊鹊结合实际经历创作
为了拯救不慎滑落下水道的一元硬币，{$t玩家}化身搜救队英勇赶赴！|鹊鹊结合实际经历创作
{$t玩家}睡着了——zzzzzzzz......|鹊鹊结合实际经历创作
在聚会上完全喝断片的{$t玩家}被半兽人三兄弟抬走咯~！♡|鹊鹊结合实际经历创作
{$t玩家}在地铁上睡着了，不断前行的车厢逐渐带他来到了最终站点...mogeko~！|鹊鹊结合实际经历创作
今天绿色章鱼俱乐部有活动，来不了了呢——by{$t玩家}|鹊鹊结合实际经历创作
“喂？跑团？啊，抱歉可能有点事情来不了了”你听着{$t玩家}电话背景音里一阵阵未知语言咏唱的声音，开始明白他现在很忙。|鹊鹊结合实际经历创作
给{$t玩家}打电话的时候，自己关注的vtb的电话也正好响了起来...|鹊鹊结合实际经历创作
因为被长发龙男逼到了小巷子，{$t玩家}大概没心思思考别的事情了。|鹊鹊结合实际经历创作
在海边散步的时候，突然被触手拉入海底的{$t玩家}！|鹊鹊结合实际经历创作
“来不了了，对不起...”电话对面的{$t玩家}房间里隐约传来阵阵喘息。|鹊鹊结合实际经历创作
黄色雨衣团真是赛高~！综上所述今天要去参加活动，来不了了哦~！——by{$t玩家}|鹊鹊结合实际经历创作
{$t玩家}正在看书，啊！不好！他被知识的巨浪冲走了！搜救队——！！！|鹊鹊结合实际经历创作
为了帮助突然晕倒的程序员木落，{$t玩家}错过了开团时间，撑住啊木落！！！|鹊鹊结合实际经历创作
由于尝试邪神召唤而来到异界的{$t玩家}，好了，这下该怎么回去呢？距离开团还有5...3...1...|鹊鹊结合实际经历创作
不慎穿越的{$t玩家}！但是接下来还有团！这一切该如何是好？《心跳！穿越到异世界了这下不得不咕咕掉跑团了呢~！》好评发售~！|鹊鹊结合实际经历创作
因为海豹一直缠着{$t玩家}，所以只好先陪他玩啦——|鹊鹊结合实际经历创作
开开心心准备开团的时候，几只大蜘蛛破窗而入！啊！{$t玩家}被他们劫走了！！！|鹊鹊结合实际经历创作
“没想到食尸鬼俱乐部的大家不是化妆特效...以后可能再也没法儿一起玩了...”{$t玩家}发来了这种意义不明的短信。|鹊鹊结合实际经历创作
“走在马路上被突如其来的龙娘威胁了，现在在小巷子里！！！请大家带一万金币救我！！！”{$t玩家}在电话里这样说。|鹊鹊结合实际经历创作
前往了一个以前捕鲸的小岛度假~这里人很亲切！但是吃了这里的鱼肉料理之后有点晕晕的诶...想到前几天{$t玩家}的短信，还是别追究他为什么不在了。|鹊鹊结合实际经历创作
因为沉迷vtb而完全忘记开团的{$t玩家}，毕竟太可爱了所以原谅他吧~！|鹊鹊结合实际经历创作
观看海豹顶球的时候站的太近被溅了一身水，换衣服的功夫{$t玩家}发现开团时间已经错过了。|鹊鹊结合实际经历创作
不知为什么平坦的路面上会躺着一只海豹，就那样玩着手机没注意就被绊倒昏过去了！可怜的{$t玩家}！|鹊鹊结合实际经历创作



{$t玩家}去依盖队大本营给大家抢香蕉了。|yumeno结合实际经历创作
“我家金鱼淹死了，要去处理一下，晚点再来”原来如此，节哀{$t玩家}！|yumeno结合实际经历创作
“我家狗在学校被老师请家长，今天不来了”这条{$t玩家}的短信让你打开手机开始搜索狗学校。|yumeno结合实际经历创作
“钱不够坐车回家，待我走回去先”{$t玩家}你其实知道手机可以支付车费的吧？|yumeno结合实际经历创作
救命！我变成鸽子了！——by{$t玩家}的短信。|yumeno结合实际经历创作
咕咕，咕咕咕咕咕，咕咕咕！——by{$t玩家}的短信。|yumeno结合实际经历创作
老板让我现在回去加班，我正在写辞呈。{$t玩家}一边内卷一边对着电话这样说。|yumeno结合实际经历创作
键盘坏了，快递还没送到，今晚不开——by{$t玩家}的短信。|yumeno结合实际经历创作
要肝活动，晚点来！——by{$t玩家}的短信。|yumeno结合实际经历创作
社区通知我现在去做核酸！by{$t玩家}的短信。|yumeno结合实际经历创作
今晚妈妈买了十斤小龙虾，可能来不了了——by{$t玩家}的短信。小龙虾是无辜的！|yumeno结合实际经历创作
“有个小孩的玩具掉轨道里了，高铁晚点了，我晚点来...是真的啦！”{$t玩家}对着手机吼道。|yumeno结合实际经历创作
“飞机没油了，我去加点油，晚点来。”——by{$t玩家}的短信。|yumeno结合实际经历创作
“寂静岭出新作了，今晚没空，咕咕咕”你看到{$t玩家}的对话框跳出这样一条内容。|yumeno结合实际经历创作
老头环中...你看着Steam好友里{$t玩家}的状态，感觉也不是不能理解。|yumeno结合实际经历创作
你打开狒狒，看见了{$t玩家}在线中，看来原因找到了。|yumeno结合实际经历创作
|yumeno结合实际经历创作


哎呀，身份证丢了，要去补办——！这条信息by{$t玩家}|秦祚轩结合实际经历创作
亲戚结婚了，我喝个喜酒就来！{$t玩家}留下这样一段话。|秦祚轩结合实际经历创作
疫情期间，主动核酸！我辈义不容辞！这样说着，{$t玩家}冲出去了。|秦祚轩结合实际经历创作
学校突然加课，大家！对不起！就算没有我你们也要从邪神手中拯救这个世界！！！{$t玩家}绝笔。|秦祚轩结合实际经历创作
滴滴滴250°c——测温枪发出这种警报，“我还有团啊——！”不理会{$t玩家}的反抗，医护人员拖走了他。|秦祚轩结合实际经历创作
钱包！我的钱包！！！不见了！！！！！！{$t玩家}一边报警一边离开了大家的视线。|秦祚轩结合实际经历创作
即使不是学雷锋日也要学雷锋！路上的老爷爷老奶奶们需要我！对不起大家！{$t玩家}在一边扶着老奶奶一边艰难的解释。|秦祚轩结合实际经历创作



你不知道今天是什么日子吗？今天是周四！你不知道周四会发生什么吗？周四有疯狂星期四！不说了，我去吃KFC了。by{$t玩家}|月森优姬结合实际经历创作



我有点事，你们先开|木落好像在结合实际经历创作
今天忽然加班了，可能来不了了|木落好像在结合实际经历创作
今天发版本，领导说发不完不让走|木落好像在结合实际经历创作
我家猫生病了，带他去看病|木落好像在结合实际经历创作
医生说今天疫苗到了，带猫打疫苗|木落好像在结合实际经历创作
我鸽某人今天就是要咕口牙！|木落好像在结合实际经历创作
当你们都觉得{$t玩家}要咕的时候，{$t玩家}咕了，这其实是一种不咕|木落好像在结合实际经历创作

{$t玩家}一觉醒来，奇怪，太阳在天上怎么还能看见星空？还有天空中这个泡泡形状的巨大黑影是什么|Szzrain结合实际经历创作

打麻将被人连胡了五个国士无双，{$t玩家}哭晕了过去——|蜜瓜包结合实际经历创作
是这样的，{$t玩家}的人格分裂被治好了，跑团的那个人格消失了，所以就完全没办法跑团啦！嗯！|蜜瓜包结合实际经历创作
什么跑团？刚分手，别来烦我！{$t玩家}如是说道|蜜瓜包结合实际经历创作
今天发大水，脑子被水淹了，跑不了团啦！|蜜瓜包结合实际经历创作
`

func RegisterBuiltinSealdiceCommands(d *Dice) {
	helpForSealdiceBlack := ".ban add user <帐号> [<原因>] //添加个人\n" +
		".ban add group <群号> [<原因>] //添加群组\n" +
		".ban add <统一ID>\n" +
		".ban rm user <帐号> //解黑/移出信任\n" +
		".ban rm group <群号>\n" +
		".ban rm <统一ID> //同上\n" +
		".ban list // 展示列表\n" +
		".ban list ban/warn/trust //只显示被禁用/被警告/信任用户\n" +
		".ban trust <统一ID> //添加信任\n" +
		".ban query <统一ID> //查看指定用户拉黑情况\n" +
		".ban help //查看帮助\n" +
		"// 统一ID示例: QQ:12345、QQ-Group:12345"
	cmdSealdiceBlack := CmdItemInfo{
		Name:      "ban",
		ShortHelp: helpForSealdiceBlack,
		Help:      "黑名单指令:\n" + helpForSealdiceBlack,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			cmdArgs.ChopPrefixToArgsWith("add", "rm", "del", "list", "show", "find", "trust")
			if ctx.PrivilegeLevel < 100 {
				ReplyToSender(ctx, msg, "你不具备Master权限")
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			getID := func() string {
				if cmdArgs.IsArgEqual(2, "user") || cmdArgs.IsArgEqual(2, "group") {
					id := cmdArgs.GetArgN(3)
					if id == "" {
						return ""
					}

					isGroup := cmdArgs.IsArgEqual(2, "group")
					return FormatDiceID(ctx, id, isGroup)
				}

				arg := cmdArgs.GetArgN(2)
				if !strings.Contains(arg, ":") {
					return ""
				}
				return arg
			}

			var val = cmdArgs.GetArgN(1)
			var uid string
			switch strings.ToLower(val) {
			case "add":
				uid = getID()
				if uid == "" {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}
				reason := cmdArgs.GetArgN(4)
				if reason == "" {
					reason = "骰主指令"
				}
				d.BanList.AddScoreBase(uid, d.BanList.ThresholdBan, "骰主指令", reason, ctx)
				ReplyToSender(ctx, msg, fmt.Sprintf("已将用户/群组 %s 加入黑名单，原因: %s", uid, reason))
			case "rm", "del":
				uid = getID()
				if uid == "" {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}

				item, ok := d.BanList.GetByID(uid)
				if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
					ReplyToSender(ctx, msg, "找不到用户/群组")
					break
				}

				ReplyToSender(ctx, msg, fmt.Sprintf("已将用户/群组 %s 移出%s列表", uid, BanRankText[item.Rank]))
				item.Score = 0
				item.Rank = BanRankNormal
			case "trust":
				uid = cmdArgs.GetArgN(2)
				if !strings.Contains(uid, ":") {
					// 如果不是这种格式，那么放弃
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}

				d.BanList.SetTrustByID(uid, "骰主指令", "骰主指令")
				ReplyToSender(ctx, msg, fmt.Sprintf("已将用户/群组 %s 加入信任列表", uid))
			case "list", "show":
				// ban/warn/trust
				var extra, text string

				extra = cmdArgs.GetArgN(2)
				d.BanList.Map.Range(func(k string, v *BanListInfoItem) bool {
					if v.Rank == BanRankNormal {
						return true
					}

					match := (extra == "trust" && v.Rank == BanRankTrusted) ||
						(extra == "ban" && v.Rank == BanRankBanned) ||
						(extra == "warn" && v.Rank == BanRankWarn)
					if extra == "" || match {
						text += v.toText(d) + "\n"
					}
					return true
				})

				if text == "" {
					text = "当前名单:\n<无内容>"
				} else {
					text = "当前名单:\n" + text
				}
				ReplyToSender(ctx, msg, text)
			case "query":
				var targetID = cmdArgs.GetArgN(2)
				if targetID == "" {
					ReplyToSender(ctx, msg, "未指定要查询的对象！")
					break
				}

				v, exists := d.BanList.Map.Load(targetID)
				if !exists {
					ReplyToSender(ctx, msg, fmt.Sprintf("所查询的<%s>情况：正常(0)", targetID))
					break
				}

				var text = fmt.Sprintf("所查询的<%s>情况：", targetID)
				switch v.Rank {
				case BanRankBanned:
					text += "禁止(-30)"
				case BanRankWarn:
					text += "警告(-10)"
				case BanRankTrusted:
					text += "信任(30)"
				default:
					text += "正常(0)"
				}
				for i, reason := range v.Reasons {
					text += fmt.Sprintf(
						"\n%s在「%s」，原因：%s",
						carbon.CreateFromTimestamp(v.Times[i]).ToDateTimeString(),
						v.Places[i],
						reason,
					)
				}
				ReplyToSender(ctx, msg, text)
			default:
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	cmdSealdiceUserID := CmdItemInfo{
		Name:      "userid",
		ShortHelp: ".userid // 查看当前帐号和群组ID",
		Help:      "查看ID:\n.userid // 查看当前帐号和群组ID",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			text := fmt.Sprintf("个人账号ID为 %s", ctx.Player.UserID)
			if !ctx.IsPrivate {
				text += fmt.Sprintf("\n当前群组ID为 %s", ctx.Group.GroupID)
			}

			ReplyToSender(ctx, msg, text)
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	helpSealdiceRoll := ".r <表达式> [<原因>] // 骰点指令\n.rh <表达式> <原因> // 暗骰"
	cmdSealdiceRoll := CmdItemInfo{
		EnableExecuteTimesParse: true,
		Name:                    "roll",
		ShortHelp:               helpSealdiceRoll,
		Help:                    "骰点:\n" + helpSealdiceRoll,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			var text string
			var diceResult int64
			var diceResultExists bool
			var detail string

			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			ctx.SystemTemplate = ctx.Group.GetCharTemplate(ctx.Dice)
			if ctx.Dice.CommandCompatibleMode {
				if (cmdArgs.Command == "rd" || cmdArgs.Command == "rhd" || cmdArgs.Command == "rdh") && len(cmdArgs.Args) >= 1 {
					if m, _ := regexp.MatchString(`^\d|优势|劣势|\+|-`, cmdArgs.CleanArgs); m {
						if cmdArgs.IsSpaceBeforeArgs {
							cmdArgs.CleanArgs = "d " + cmdArgs.CleanArgs
						} else {
							cmdArgs.CleanArgs = "d" + cmdArgs.CleanArgs
						}
					}
				}
			}

			var r *VMResultV2m
			var commandInfoItems []any

			rollOne := func() *CmdExecuteResult {
				forWhat := ""
				var matched string

				if len(cmdArgs.Args) >= 1 { //nolint:nestif
					var err error
					r, detail, err = DiceExprEvalBase(ctx, cmdArgs.CleanArgs, RollExtraFlags{
						DefaultDiceSideNum: getDefaultDicePoints(ctx),
						DisableBlock:       true,
						V2Only:             true,
					})

					if r != nil && !r.IsCalculated() {
						forWhat = cmdArgs.CleanArgs

						defExpr := "d"
						if ctx.diceExprOverwrite != "" {
							defExpr = ctx.diceExprOverwrite
						}
						r, detail, err = DiceExprEvalBase(ctx, defExpr, RollExtraFlags{
							DefaultDiceSideNum: getDefaultDicePoints(ctx),
							DisableBlock:       true,
						})
					}

					if r != nil && r.TypeId == ds.VMTypeInt {
						diceResult = int64(r.MustReadInt())
						diceResultExists = true
					}

					if err == nil {
						matched = r.GetMatched()
						if forWhat == "" {
							forWhat = r.GetRestInput()
						}
					} else {
						errs := err.Error()
						if strings.HasPrefix(errs, "E1:") || strings.HasPrefix(errs, "E5:") || strings.HasPrefix(errs, "E6:") || strings.HasPrefix(errs, "E7:") || strings.HasPrefix(errs, "E8:") {
							ReplyToSender(ctx, msg, errs)
							return &CmdExecuteResult{Matched: true, Solved: true}
						}
						forWhat = cmdArgs.CleanArgs
					}
				}

				VarSetValueStr(ctx, "$t原因", forWhat)
				if forWhat != "" {
					forWhatText := DiceFormatTmpl(ctx, "核心:骰点_原因")
					VarSetValueStr(ctx, "$t原因句子", forWhatText)
				} else {
					VarSetValueStr(ctx, "$t原因句子", "")
				}

				if diceResultExists { //nolint:nestif
					detailWrap := ""
					if detail != "" {
						detailWrap = "=" + detail
						re := regexp.MustCompile(`\[((\d+)d\d+)\=(\d+)\]`)
						match := re.FindStringSubmatch(detail)
						if len(match) > 0 {
							num := match[2]
							if num == "1" && (match[1] == matched || match[1] == "1"+matched) {
								detailWrap = ""
							}
						}
					}

					// 指令信息标记
					item := map[string]interface{}{
						"expr":   matched,
						"result": diceResult,
						"reason": forWhat,
					}
					if forWhat == "" {
						delete(item, "reason")
					}
					commandInfoItems = append(commandInfoItems, item)

					VarSetValueStr(ctx, "$t表达式文本", matched)
					VarSetValueStr(ctx, "$t计算过程", detailWrap)
					VarSetValueInt64(ctx, "$t计算结果", diceResult)
				} else {
					var val int64
					var detail string
					dicePoints := getDefaultDicePoints(ctx)
					if ctx.diceExprOverwrite != "" {
						r, detail, _ = DiceExprEvalBase(ctx, cmdArgs.CleanArgs, RollExtraFlags{
							DefaultDiceSideNum: dicePoints,
							DisableBlock:       true,
						})
						if r != nil && r.TypeId == ds.VMTypeInt {
							valX, _ := r.ReadInt()
							val = int64(valX)
						}
					} else {
						r, _, _ = DiceExprEvalBase(ctx, "d", RollExtraFlags{
							DefaultDiceSideNum: dicePoints,
							DisableBlock:       true,
						})
						if r != nil && r.TypeId == ds.VMTypeInt {
							valX, _ := r.ReadInt()
							val = int64(valX)
						}
					}

					// 指令信息标记
					item := map[string]any{
						"expr":       fmt.Sprintf("D%d", dicePoints),
						"reason":     forWhat,
						"dicePoints": dicePoints,
						"result":     val,
					}
					if forWhat == "" {
						delete(item, "reason")
					}
					commandInfoItems = append(commandInfoItems, item)

					VarSetValueStr(ctx, "$t表达式文本", fmt.Sprintf("D%d", dicePoints))
					VarSetValueStr(ctx, "$t计算过程", detail)
					VarSetValueInt64(ctx, "$t计算结果", val)
				}
				return nil
			}

			if cmdArgs.SpecialExecuteTimes > 1 {
				VarSetValueInt64(ctx, "$t次数", int64(cmdArgs.SpecialExecuteTimes))
				if cmdArgs.SpecialExecuteTimes > int(ctx.Dice.MaxExecuteTime) {
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰点_轮数过多警告"))
					return CmdExecuteResult{Matched: true, Solved: true}
				}
				var texts []string
				for i := 0; i < cmdArgs.SpecialExecuteTimes; i++ {
					ret := rollOne()
					if ret != nil {
						return *ret
					}
					texts = append(texts, DiceFormatTmpl(ctx, "核心:骰点_单项结果文本"))
				}
				VarSetValueStr(ctx, "$t结果文本", strings.Join(texts, "\n"))
				text = DiceFormatTmpl(ctx, "核心:骰点_多轮")
			} else {
				ret := rollOne()
				if ret != nil {
					return *ret
				}
				VarSetValueStr(ctx, "$t结果文本", DiceFormatTmpl(ctx, "核心:骰点_单项结果文本"))
				text = DiceFormatTmpl(ctx, "核心:骰点")
			}

			isHide := strings.Contains(cmdArgs.Command, "h")

			// 指令信息
			commandInfo := map[string]any{
				"cmd":    "roll",
				"pcName": ctx.Player.Name,
				"items":  commandInfoItems,
			}
			if isHide {
				commandInfo["hide"] = isHide
			}
			ctx.CommandInfo = commandInfo

			if kw := cmdArgs.GetKwarg("asm"); r != nil && kw != nil {
				if ctx.PrivilegeLevel >= 40 {
					asm := r.GetAsmText()
					text += "\n" + asm
				}
			}

			if kw := cmdArgs.GetKwarg("ci"); kw != nil {
				info, err := json.Marshal(ctx.CommandInfo)
				if err == nil {
					text += "\n" + string(info)
				} else {
					text += "\n" + "指令信息无法序列化"
				}
			}

			if isHide {
				if msg.Platform == "QQ-CH" {
					ReplyToSender(ctx, msg, "QQ频道内尚不支持暗骰")
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				if ctx.Group != nil {
					if ctx.IsPrivate {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_私聊不可用"))
					} else {
						ctx.CommandHideFlag = ctx.Group.GroupID
						prefix := DiceFormatTmpl(ctx, "核心:暗骰_私聊_前缀")
						ReplyGroup(ctx, msg, DiceFormatTmpl(ctx, "核心:暗骰_群内"))
						ReplyPerson(ctx, msg, prefix+text)
					}
				} else {
					ReplyToSender(ctx, msg, text)
				}
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			ReplyToSender(ctx, msg, text)
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	helpSealdiceRollX := ".rx <表达式> <原因> // 骰点指令\n.rxh <表达式> <原因> // 暗骰"
	cmdSealdiceRollX := CmdItemInfo{
		Name:          "roll",
		ShortHelp:     helpSealdiceRollX,
		Help:          "骰点(和r相同，但支持代骰):\n" + helpSealdiceRollX,
		AllowDelegate: true,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			mctx := GetCtxProxyFirst(ctx, cmdArgs)
			return cmdSealdiceRoll.Solve(mctx, msg, cmdArgs)
		},
	}

	cmdRsr := CmdItemInfo{
		Name:      "rsr",
		ShortHelp: ".rsr <骰数> // 暗影狂奔",
		Help: "暗影狂奔骰点:\n.rsr <骰数>\n" +
			"> 每个被骰出的五或六就称之为一个成功度\n" +
			"> 如果超过半数的骰子投出了一被称之为失误\n" +
			"> 在投出失误的同时没能骰出至少一个成功度被称之为严重失误",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			val := cmdArgs.GetArgN(1)
			num, err := strconv.ParseInt(val, 10, 64)

			if err == nil && num > 0 {
				successDegrees := int64(0)
				failedCount := int64(0)
				var results []string
				for i := int64(0); i < num; i++ {
					v := DiceRoll64(6)
					if v >= 5 {
						successDegrees++
					} else if v == 1 {
						failedCount++
					}
					// 过大的骰池不显示
					if num < 10 {
						results = append(results, strconv.FormatInt(v, 10))
					}
				}

				var detail string
				if len(results) > 0 {
					detail = "{" + strings.Join(results, "+") + "}\n"
				}

				text := fmt.Sprintf("<%s>骰点%dD6:\n", ctx.Player.Name, num)
				text += detail
				text += fmt.Sprintf("成功度:%d/%d\n", successDegrees, failedCount)

				successRank := int64(0) // 默认
				if failedCount > (num / 2) {
					// 半数失误
					successRank = -1

					if successDegrees == 0 {
						successRank = -2
					}
				}

				switch successRank {
				case -1:
					text += "失误"
				case -2:
					text += "严重失误"
				}
				ReplyToSender(ctx, msg, text)
			} else {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	// Emoklore(共鸣性怪异)规则支持
	helpEk := ".ek <技能名称>[+<奖励骰>] 判定值\n" +
		".ek 检索 // 骰“检索”等级个d10，计算成功数\n" +
		".ek 检索+2 // 在上一条基础上加骰2个d10\n" +
		".ek 检索 6  // 骰“检索”等级个d10，计算小于6的骰个数\n" +
		".ek 检索 知力+检索 // 骰”检索“，判定线为”知力+检索“\n" +
		".ek 5 4 // 骰5个d10，判定值4\n" +
		".ek 检索2 // 未录卡情况下判定2级检索\n" +
		".ek 共鸣 6 // 共鸣判定，成功后手动st共鸣+N\n"
	cmdEk := CmdItemInfo{
		Name:      "ek",
		ShortHelp: helpEk,
		Help:      "共鸣性怪异骰点:\n" + helpEk,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			mctx := ctx

			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			txt := cmdArgs.CleanArgs
			re := regexp.MustCompile(`(?:([^*+\-\s\d]+)(\d+)?|(\d+))\s*(?:([+\-*])\s*(\d+))?`)
			m := re.FindStringSubmatch(txt)
			if len(m) == 0 {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			// 读取技能名字和等级
			mustHaveCheckVal := false
			name := m[1]         // .ek 摸鱼
			nameLevelStr := m[2] // .ek 摸鱼3
			if name == "" && nameLevelStr == "" {
				// .ek 3 4
				nameLevelStr = m[3]
				mustHaveCheckVal = true
			}

			var nameLevel int64
			if nameLevelStr != "" {
				nameLevel, _ = strconv.ParseInt(nameLevelStr, 10, 64)
			} else {
				nameLevel, _ = VarGetValueInt64(mctx, name)
			}

			// 附加值 .ek 技能+1
			extraOp := m[4]
			extraValStr := m[5]
			extraVal := int64(0)
			if extraValStr != "" {
				extraVal, _ = strconv.ParseInt(extraValStr, 10, 64)
				if extraOp == "-" {
					extraVal = -extraVal
				}
			}

			restText := txt[len(m[0]):]
			restText = strings.TrimSpace(restText)

			if restText == "" && mustHaveCheckVal {
				ReplyToSender(ctx, msg, "必须填入判定值")
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			// 填充补充部分
			if restText == "" {
				restText = fmt.Sprintf("%s%s", name, nameLevelStr)
				mode := 1
				v := emokloreAttrParent[name]
				if v == nil {
					v = emokloreAttrParent2[name]
					mode = 2
				}
				if v == nil {
					v = emokloreAttrParent3[name]
					mode = 3
				}
				if v != nil {
					maxName := ""
					maxVal := int64(0)
					for _, i := range v {
						val, _ := VarGetValueInt64(mctx, i)
						if val >= maxVal {
							maxVal = val
							maxName = i
						}
					}
					if maxName != "" {
						switch mode {
						case 1:
							// 种类1: 技能+属性
							restText += " + " + maxName
						case 2:
							// 种类2: 属性/2[向上取整]
							restText = fmt.Sprintf("(%s+1)/2", maxName)
						case 3:
							// 种类3: 属性
							restText = maxName
						}
					}
				}
			}

			r, detail, err := mctx.Dice._ExprEvalBaseV1(restText, mctx, RollExtraFlags{
				CocVarNumberMode: true,
				DisableBlock:     true,
			})
			if err != nil {
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			checkVal, _ := r.ReadInt64()
			diceNum := nameLevel // 骰子个数为技能等级，至少1个
			if diceNum < 1 {
				diceNum = 1
			}
			if extraOp == "*" {
				diceNum *= extraVal
			} else {
				diceNum += extraVal
			}

			successDegrees := int64(0)
			var results []string
			for i := int64(0); i < diceNum; i++ {
				v := DiceRoll64(10)
				if v <= checkVal {
					successDegrees++
				}
				if v == 1 {
					successDegrees++
				}
				if v == 10 {
					successDegrees--
				}
				// 过大的骰池不显示
				if diceNum < 15 {
					results = append(results, strconv.FormatInt(v, 10))
				}
			}

			var detailPool string
			if len(results) > 0 {
				detailPool = "{" + strings.Join(results, "+") + "}\n"
			}

			// 检定原因
			showName := name
			if showName == "" {
				showName = nameLevelStr
			}
			if nameLevelStr != "" {
				showName += nameLevelStr
			}
			if extraVal > 0 {
				showName += extraOp + extraValStr
			}

			if detail != "" {
				detail = "{" + detail + "}"
			}

			checkText := ""
			switch {
			case successDegrees < 0:
				checkText = "大失败"
			case successDegrees == 0:
				checkText = "失败"
			case successDegrees == 1:
				checkText = "通常成功"
			case successDegrees == 2:
				checkText = "有效成功"
			case successDegrees == 3:
				checkText = "极限成功"
			case successDegrees >= 10:
				checkText = "灾难成功"
			case successDegrees >= 4:
				checkText = "奇迹成功"
			}

			text := fmt.Sprintf("<%s>的“%s”共鸣性怪异规则检定:\n", ctx.Player.Name, showName)
			text += detailPool
			text += fmt.Sprintf("判定值: %d%s\n", checkVal, detail)
			text += fmt.Sprintf("成功数: %d[%s]\n", successDegrees, checkText)

			ReplyToSender(ctx, msg, text)
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	helpEkGen := ".ekgen [<数量>] // 制卡指令，生成<数量>组人物属性，最高为10次"
	cmdEkgen := CmdItemInfo{
		Name:      "ekgen",
		ShortHelp: helpEkGen,
		Help:      "共鸣性怪异制卡指令:\n" + helpEkGen,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			n := cmdArgs.GetArgN(1)
			val, err := strconv.ParseInt(n, 10, 64)
			if err != nil {
				if n == "" {
					val = 1 // 数量不存在时，视为1次
				} else {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}
			}
			if val > 10 {
				val = 10
			}
			var i int64

			var ss []string
			for i = 0; i < val; i++ {
				randMap := map[int64]bool{}
				for j := 0; j < 6; j++ {
					n := DiceRoll64(24)
					if randMap[n] {
						j-- // 如果已经存在，重新roll
					} else {
						randMap[n] = true
					}
				}

				var nums Int64SliceDesc
				for k := range randMap {
					nums = append(nums, k)
				}
				sort.Sort(nums)

				last := int64(25)
				var nums2 []interface{}
				for _, j := range nums {
					val := last - j
					last = j
					nums2 = append(nums2, val)
				}
				nums2 = append(nums2, last)

				// 过滤大于6的
				for {
					// 遍历找出一个大于6的
					isGT6 := false
					var rest int64
					for index, _j := range nums2 {
						j := _j.(int64)
						if j > 6 {
							isGT6 = true
							rest = j - 6
							nums2[index] = int64(6)
							break
						}
					}

					if isGT6 {
						for index, _j := range nums2 {
							j := _j.(int64)
							if j < 6 {
								nums2[index] = j + rest
								break
							}
						}
					} else {
						break
					}
				}
				rand.Shuffle(len(nums2), func(i, j int) {
					nums2[i], nums2[j] = nums2[j], nums2[i]
				})

				text := fmt.Sprintf("身体:%d 灵巧:%d 精神:%d 五感:%d 知力:%d 魅力:%d 社会:%d", nums2...)
				text += fmt.Sprintf(" 运势:%d hp:%d mp:%d", DiceRoll64(6), nums2[0].(int64)+10, nums2[2].(int64)+nums2[4].(int64))

				ss = append(ss, text)
			}
			info := strings.Join(ss, "\n")
			ReplyToSender(ctx, msg, fmt.Sprintf("<%s>的共鸣性怪异人物做成:\n%s", ctx.Player.Name, info))
			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	readNumber := func(text string, extra string) string {
		if text == "" {
			return ""
		}
		re0 := regexp.MustCompile(`^((\d+)[dDcCaA]|[bBpPfF])(.*)`)
		if re0.MatchString(text) {
			// 这种不需要管，是合法的表达式
			return text
		}

		re := regexp.MustCompile(`^(\d+)(.*)`)
		m := re.FindStringSubmatch(text)
		if len(m) > 0 {
			var rest string
			if len(m) > 2 {
				rest = m[2]
			}
			// 数字 a10 剩下部分
			return fmt.Sprintf("%s%s%s", m[1], extra, rest)
		}

		return text
	}

	cmdDX := CmdItemInfo{
		Name:      "dx",
		ShortHelp: ".dx 3c4",
		Help:      "双重十字规则骰点:\n.dx 3c4 // 也可使用.r 3c4替代",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			txt := readNumber(cmdArgs.CleanArgs, "c10")
			if txt == "" {
				txt = "1c10"
				cmdArgs.Args = []string{txt}
			}
			cmdArgs.CleanArgs = txt
			ctx.diceExprOverwrite = "1c10"
			roll := ctx.Dice.CmdMap["roll"]
			return roll.Solve(ctx, msg, cmdArgs)
		},
	}

	cmdJsr := CmdItemInfo{
		EnableExecuteTimesParse: true,
		Name:                    "jsr",
		ShortHelp:               ".jsr 3# 10 // 投掷 10 面骰 3 次，结果不重复。结果存入骰池并可用 .drl 抽取。",
		Help: "不重复骰点(Jetter sans répéter):\n.jsr <次数># <投骰表达式> [<名字>]" +
			"\n用例：.jsr 3# 10 // 投掷 10 面骰 3 次，结果不重复，结果存入骰池并可用 .drl 抽取。",
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			t := cmdArgs.SpecialExecuteTimes
			allArgClean := cmdArgs.CleanArgs
			allArgs := strings.Split(allArgClean, " ")
			var m int
			for i, v := range allArgs {
				if strings.HasPrefix(v, "d") {
					v = strings.Replace(v, "d", "", 1)
				}

				if n, err := strconv.Atoi(v); err == nil {
					m = n
					allArgs = append(allArgs[:i], allArgs[i+1:]...)
					break
				}
			}
			if t == 0 {
				t = 1
			}
			if m == 0 {
				m = int(getDefaultDicePoints(ctx))
			}
			if t > int(ctx.Dice.MaxExecuteTime) {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰点_轮数过多警告"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			if t > m {
				ReplyToSender(ctx, msg, fmt.Sprintf("无法不重复地投掷%d次%d面骰。", t, m))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			var pool []int
			ma := make(map[int]bool)
			for len(pool) < t {
				n := rand.IntN(m) + 1
				if !ma[n] {
					ma[n] = true
					pool = append(pool, n)
				}
			}
			var results []string
			for _, v := range pool {
				results = append(results, fmt.Sprintf("D%d=%d", m, v))
			}
			allArgClean = strings.Join(allArgs, " ")
			for i := range pool {
				j := rand.IntN(i + 1)
				pool[i], pool[j] = pool[j], pool[i]
			}
			roulette := singleRoulette{
				Face: int64(m),
				Name: allArgClean,
				Pool: pool,
			}

			rouletteMap.Store(ctx.Group.GroupID, roulette)
			VarSetValueStr(ctx, "$t原因", allArgClean)
			if allArgClean != "" {
				forWhatText := DiceFormatTmpl(ctx, "核心:骰点_原因")
				VarSetValueStr(ctx, "$t原因句子", forWhatText)
			} else {
				VarSetValueStr(ctx, "$t原因句子", "")
			}
			VarSetValueInt64(ctx, "$t次数", int64(t))
			VarSetValueStr(ctx, "$t结果文本", strings.Join(results, "\n"))
			reply := DiceFormatTmpl(ctx, "核心:骰点_多轮")
			ReplyToSender(ctx, msg, reply)
			return CmdExecuteResult{
				Matched: true,
				Solved:  true,
			}
		},
	}

	cmdDrl := CmdItemInfo{
		EnableExecuteTimesParse: true,
		Name:                    "drl",
		ShortHelp: ".drl new 10 5# // 在当前群组创建一个面数为 10，能抽取 5 次的骰池\n.drl // 抽取当前群组的骰池\n" +
			".drlh //抽取当前群组的骰池，结果私聊发送",
		Help: "drl（Draw Lot）：.drl new <次数> <投骰表达式> [<名字>] // 在当前群组创建一个骰池\n" +
			"用例：.drl new 10 5# // 在当前群组创建一个面数为 10，能抽取 5 次的骰池\n\n.drl // 抽取当前群组的骰池\n" +
			".drlh //抽取当前群组的骰池，结果私聊发送",
		DisabledInPrivate: true,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			if cmdArgs.IsArgEqual(1, "help") {
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
			if cmdArgs.IsArgEqual(1, "new") {
				// Make mode
				roulette := singleRoulette{
					Name: "",
					Face: getDefaultDicePoints(ctx),
					Time: 1,
				}
				t := cmdArgs.SpecialExecuteTimes
				if t != 0 {
					roulette.Time = t
				}

				m := cmdArgs.GetArgN(2)
				n := m
				if strings.HasPrefix(m, "d") {
					m = strings.Replace(m, "d", "", 1)
				}
				if i, err := strconv.Atoi(m); err == nil {
					roulette.Face = int64(i)
					text := cmdArgs.GetArgN(3)
					roulette.Name = text
				} else {
					roulette.Name = n
				}

				// NOTE(Xiangze Li): 允许创建更多轮数。使用洗牌算法后并不会很重复计算
				// if roulette.Time > int(ctx.Dice.MaxExecuteTime) {
				// 	ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰点_轮数过多警告"))
				// 	return CmdExecuteResult{Matched: true, Solved: true}
				// }

				if int64(roulette.Time) > roulette.Face {
					ReplyToSender(ctx, msg, fmt.Sprintf("创建错误：无法不重复地投掷%d次%d面骰。",
						roulette.Time,
						roulette.Face))
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				// 创建pool后产生随机数，使用F-Y洗牌算法以保证随机性和效率
				var pool = make([]int, roulette.Time)
				var allNum = make([]int, roulette.Face)
				for i := range allNum {
					allNum[i] = i + 1
				}
				for idx := 0; idx < roulette.Time; idx++ {
					i := int(roulette.Face) - 1 - idx
					j := rand.IntN(i + 1)
					allNum[i], allNum[j] = allNum[j], allNum[i]
					pool[idx] = allNum[i]
				}
				roulette.Pool = pool

				rouletteMap.Store(ctx.Group.GroupID, roulette)
				ReplyToSender(ctx, msg, fmt.Sprintf("创建骰池%s成功，骰子面数%d，可抽取%d次。",
					roulette.Name, roulette.Face, roulette.Time))
				return CmdExecuteResult{Matched: true, Solved: true}
			}
			// Draw mode
			var isRouletteEmpty = true
			rouletteMap.Range(func(key string, value singleRoulette) bool {
				isRouletteEmpty = false
				return false
			})
			tryLoad, ok := rouletteMap.Load(ctx.Group.GroupID)
			if isRouletteEmpty || !ok || tryLoad.Face == 0 {
				ReplyToSender(ctx, msg, "当前群组无骰池，请使用.drl new创建一个。")
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			result := fmt.Sprintf("D%d=%d", tryLoad.Face, tryLoad.Pool[0])
			tryLoad.Pool = append(tryLoad.Pool[:0], tryLoad.Pool[1:]...)
			VarSetValueStr(ctx, "$t原因", tryLoad.Name)
			if tryLoad.Name != "" {
				forWhatText := DiceFormatTmpl(ctx, "核心:骰点_原因")
				VarSetValueStr(ctx, "$t原因句子", forWhatText)
			} else {
				VarSetValueStr(ctx, "$t原因句子", "")
			}
			VarSetValueStr(ctx, "$t结果文本", result)
			reply := DiceFormatTmpl(ctx, "核心:骰点")

			if cmdArgs.Command == "drl" {
				if len(tryLoad.Pool) == 0 {
					reply += "\n骰池已经抽空，现在关闭。"
					tryLoad = singleRoulette{}
				}
				ReplyToSender(ctx, msg, reply)
			} else if cmdArgs.Command == "drlh" {
				announce := msg.Sender.Nickname + "进行了抽取。"
				reply += fmt.Sprintf("\n来自群%s(%s)",
					ctx.Group.GroupName, ctx.Group.GroupID)
				if len(tryLoad.Pool) == 0 {
					announce += "\n骰池已经抽空，现在关闭。"
					tryLoad = singleRoulette{}
				}
				ReplyGroup(ctx, msg, announce)
				ReplyPerson(ctx, msg, reply)
			}
			rouletteMap.Store(ctx.Group.GroupID, tryLoad)
			return CmdExecuteResult{
				Matched: true,
				Solved:  true,
			}
		},
	}

	cmdSealBot := CmdItemInfo{
		Name:      "SealdiceBot",
		ShortHelp: ".bot on/off/about/bye/quit // 开启、关闭、查看信息、退群",
		Help:      "骰子管理:\n.bot on/off/about/bye[exit,quit] // 开启、关闭、查看信息、退群",
		Raw:       true,
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

			if len(cmdArgs.Args) > 0 && !cmdArgs.IsArgEqual(1, "about") { //nolint:nestif
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

				if cmdArgs.IsArgEqual(1, "on") {
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
					if ctx.Group.LogOn {
						text += "\n请特别注意: 日志记录处于开启状态"
					}
					ReplyToSender(ctx, msg, text)

					return CmdExecuteResult{Matched: true, Solved: true}
				} else if cmdArgs.IsArgEqual(1, "off") {
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
				} else if cmdArgs.IsArgEqual(1, "bye", "exit", "quit") {
					if cmdArgs.GetArgN(2) != "" {
						return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
					}

					if ctx.IsPrivate {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_私聊不可用"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					if ctx.PrivilegeLevel < 40 {
						ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master/管理"))
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					if !cmdArgs.AmIBeMentioned {
						// 裸指令，如果当前群内开启，予以提示
						if ctx.IsCurGroupBotOn {
							ReplyToSender(ctx, msg, "[退群指令] 请@我使用这个命令，以进行确认")
						}
						return CmdExecuteResult{Matched: true, Solved: true}
					}

					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰子退群预告"))

					userName := ctx.Dice.Parent.TryGetUserName(msg.Sender.UserID)
					txt := fmt.Sprintf("指令退群: 于群组<%s>(%s)中告别，操作者:<%s>(%s)",
						ctx.Group.GroupName, msg.GroupID, userName, msg.Sender.UserID)
					d.Logger.Info(txt)
					ctx.Notice(txt)

					// SetBotOffAtGroup(ctx, ctx.Group.GroupID)
					time.Sleep(3 * time.Second)
					ctx.Group.DiceIDExistsMap.Delete(ctx.EndPoint.UserID)
					ctx.Group.UpdatedAtTime = time.Now().Unix()
					ctx.EndPoint.Adapter.QuitGroup(ctx, msg.GroupID)

					return CmdExecuteResult{Matched: true, Solved: true}
				} else if cmdArgs.IsArgEqual(1, "save") {
					d.Save(false)
					ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:骰子保存设置"))
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}

			if cmdArgs.SomeoneBeMentionedButNotMe {
				return CmdExecuteResult{Matched: false, Solved: false}
			}

			activeCount := 0
			serveCount := 0
			// Pinenutn: Range模板 ServiceAtNew重构代码
			d.ImSession.ServiceAtNew.Range(func(_ string, gp *GroupInfo) bool {
				// Pinenutn: ServiceAtNew重构
				if gp.GroupID != "" &&
					!strings.HasPrefix(gp.GroupID, "PG-") &&
					gp.DiceIDExistsMap.Exists(ctx.EndPoint.UserID) {
					serveCount++
					if gp.DiceIDActiveMap.Exists(ctx.EndPoint.UserID) {
						activeCount++
					}
				}
				return true
			})

			onlineVer := ""
			if d.Parent.AppVersionOnline != nil {
				ver := d.Parent.AppVersionOnline
				// 如果当前不是最新版，那么提示
				if ver.VersionLatestCode != VERSION_CODE {
					onlineVer = "\n最新版本: " + ver.VersionLatestDetail + "\n"
				}
			}

			var groupWorkInfo, activeText string
			if inGroup {
				activeText = "关闭"
				if ctx.Group.IsActive(ctx) {
					activeText = "开启"
				}
				groupWorkInfo = "\n群内工作状态: " + activeText
			}

			VarSetValueInt64(ctx, "$t供职群数", int64(serveCount))
			VarSetValueInt64(ctx, "$t启用群数", int64(activeCount))
			VarSetValueStr(ctx, "$t群内工作状态", groupWorkInfo)
			VarSetValueStr(ctx, "$t群内工作状态_仅状态", activeText)
			ver := VERSION.String()
			arch := runtime.GOARCH
			if arch != "386" && arch != "amd64" {
				ver = fmt.Sprintf("%s %s", ver, arch)
			}
			baseText := fmt.Sprintf("SealDice %s%s", ver, onlineVer)
			extText := DiceFormatTmpl(ctx, "核心:骰子状态附加文本")
			if extText != "" {
				extText = "\n" + extText
			}
			text := baseText + extText

			ReplyToSender(ctx, msg, text)

			return CmdExecuteResult{Matched: true, Solved: true}
		},
	}

	d.RegisterExtension(&ExtInfo{
		Name:            "SealdiceCore", // 扩展的名称，需要用于指令中，写简短点      2024.05.10: 目前被看成是 function 的缩写了（
		Version:         "1.0.0",
		Brief:           "海豹核心指令",
		AutoActive:      false, // 是否自动开启
		ActiveOnPrivate: true,
		Author:          "海棠",
		Official:        true,
		OnCommandReceived: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) {
		},
		OnLoad: func() {
		},
		GetDescText: GetExtensionDesc,
		CmdMap: CmdMapCls{
			"ban":    &cmdSealdiceBlack,
			"black":  &cmdSealdiceBlack,
			"userid": &cmdSealdiceUserID,
			"rd":     &cmdSealdiceRoll,
			"roll":   &cmdSealdiceRoll,
			"rhd":    &cmdSealdiceRoll,
			"rdh":    &cmdSealdiceRoll,
			"rx":     &cmdSealdiceRollX,
			"rxh":    &cmdSealdiceRollX,
			"rhx":    &cmdSealdiceRollX,
			"rsr":    &cmdRsr,
			"ek":     &cmdEk,
			"ekgen":  &cmdEkgen,
			"dx":     &cmdDX,
			"dxh":    &cmdDX,
			"jsr":    &cmdJsr,
			"drl":    &cmdDrl,
			"drlh":   &cmdDrl,
			"bot":    &cmdSealBot,
		},
	})
}
