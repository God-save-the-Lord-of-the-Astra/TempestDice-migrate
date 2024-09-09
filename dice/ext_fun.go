package dice

import (
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"

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
var duelText = `
完胜，这是彻彻底底的完胜啊哈哈哈
啧，就差一点就完胜了（懊恼）
看看！看看！这就是自不量力来挑战我的下场
哼哼，看这天堑般的差距！
认输了没有，嘿嘿，见识到实力的差距了吧
木大木大，至少再回去修炼个几百年爆裂魔法再来挑战我
(得意洋洋的目光)
你要是现在向我求饶我也不是不能原谅你哦.嗯？还敢不敢挑战我了？（笑）
明天再来也是一样的下场嘿嘿（叉腰）
木大木大木大，全部木大
这么悬殊的比分可是少见的哦？
我觉得你可以跟悠悠学学
还要来吗？再来也是一样的啦
要不要我下次放你一马？
嘿嘿嘿，我可绝对没作弊
嘛，既然输了就付出代价吧，（拽住）来来，我看那边那家店的汉堡肉不错
嘛，既然输了就付出代价吧，看见那边正在嚼嚼嚼的miku了没，把她手里的大葱抢来给我
看起来你很不服嘛,但是不服也没用，我赢了就是我赢了
嘛既然。。。既然输了就要付出代价，那个（扭捏），嗯，给我一个公主抱试试看
呼呼，是我的完胜呢，认输了没有？
你觉得靠次数就能赢我那你就错了嘿嘿
嘛，又输了，也很正常嘛，毕竟！对手是我
嘛，你又输了（揉揉脑袋）
嘛，我赢了，我要吃薯片（趴在沙发上）
表面上这是一个概率问题，实际上，就是我很强啦
啧 啧 啧，太弱了，换个人来吧
你！输！了！不服？可以啊，明天再来，我再赢你一次（笑）
再试也没用啦，会输一次就会输一百次哦
嗯哼~是我的大胜利
不服？不服你就正面赢我啊
喂。。。你盯着我干什么。。没有！我说没有灌铅就没有灌铅，这是公平的决斗。。。嗯。。公平的
你觉得能赢我，嘿嘿，那是个幻觉，美丽的幻觉。
之前也有个调查员像你那么自信，后来他被撕成了两半
哦吼？我觉得我甚至可以让你一个骰子
（拍肩）不要那么沮丧嘛，人生不就是大起大落落落落落落落落
你请我吃顿饭这次就算你赢，怎么样？
我觉得吧，你再来一百次，也是这个结果（笑）
安啦，是我太厉害不是你太菜
嗯哼？比分居然比我预料的还要高一点，算你有两下子
(哼起了歌)
来来来既然输了就不要不服，过来让我揍一拳
哦吼，来来来就像说好的那样，输的人去给682换洗澡水
（面露笑意地看着你）
你输了~你输了~你输了！重要的事说三遍
我可是骰娘哦？你怎么可能赢嘛
赢是不可能让你赢的哦
啊，居然被你拉到了这个分数，那就夸你一下好了
看似你只差一分就赢了，实则这一分是天差地别
这不是平局！我赢了！我说我赢了就是我赢了！
库唔，不可能。。居然。这是我算错了，肯定是我算错了。
居然。。。不对。。我赢了。。一定是我赢了（念叨）
你觉得是你赢了就错了！这是Error，肯定是哪里出错了！
？！居然趁我一时大意。。。
这不对。。哪里不对！我不相信！
你你你。。。我警告你，不要得意，我我明天就找回场子
玛丽安姐！！！！！有人欺负我！打他！
假。。。假的吧（跌倒）虽然只有几分但是居然。。让你赢了
木大木大，都说了我让你一只手也。。。嗯？嗯？！（长久的沉默）
这是幻觉。。。这一定是幻觉。。
呜呜呜。你欺负我，我才没有哭！你给我记着！（跑开）
（盯————）我记住你了
（捂住头）什么嘛！明明灌了铅为什么还是输了
我要求重来！重来！这个不算数！
在我心里我赢了！分数不重要！嗯，实际分数没有任何意义！
（小声）你等着，我一定是要用大失败找回场子
这是什么情况，你是不是作弊了（嚼嚼嚼）
你肯定作弊了对不对？
哼！这次就算你赢了吧，愿赌服输，你想要啥？
那啥。。。咱们当做什么都没发生好不好（递水）
......@惠惠.botoff
（看起来灌的铅还不够。。）咳咳，什么都没有，这次就算你运气好这次就算你运气好，哼，我看上去像是那么死不认账的人吗？
嗯，我赢了，你不要说话，对，我赢了，我赢了，我我我我说我赢了就是我赢了！！我才没有哭！
只不过是你一时运气好了那么一点点罢了
哼，不要得意，等明天。。。
这什么啊！还是我终于开始出现幻觉了？！为什么看起来好像是你赢了
唔，哼，虽然看起来是你赢了，但是但是但是（声音逐渐变小）
呐，你赢了，呐呐，按传统，赢的人要请吃饭啊
这。。。（陷入沉思）
我。。我居然输了，你你你你你
唉？？巴尼尔先生明明说我今天会赢的啊
（掀桌）
能。。再来一次吗，我我我，我让你见识一下我的爆裂魔法作为交换
这...嗯，这也在计算之中，我早就算到了你的分数会高一点，所以其实还是我赢了
呐，你赢了，哼哼，我可没有那么小家子气，我可是红魔族第一的天才啊
（移开视线默默遮住骰子）你，，你什么都没看见，我什么都没有扔
唔，大家，朋友一场，这个这个，就当做我赢了怎么样啊，我我我我请你吃小龙虾，我的小龙虾料理可是一流的啊
唔，今天我看起来状态不大好呢，因为你看，这个数字，原本那里应该有个负号的吧，为什么我看不清？
肯定是因为刚刚放完爆裂魔法状态不好，是的，是这样的
....敢不敢再来一次啊kura
(长久的沉默)
你看。。看在我给你投了这么多次骰子的份上。。大家当做什么都没有发生如何
我觉得我们肯定有什么误会在哒（搓手），不如大家交个朋友怎么样
？！居然。。如此悬殊的比分。。不可能不可能，刚刚程序出bug了，不如我们再来一次（小声）
(摇摇摇)这个骰子质量有问题，不如我们换一个再来一次？
你你你你，你欺负我，我要告诉玛丽安姐，回头她会来收拾你的！你给我记住！记住！
！！！！！！！（一脸不可置信地原地晕倒）
?居然。。我居然会输到这种地步。。你你你，你欺负我，我要退群！退群！
完败？！这不可能？不可能？！你你。。你对我的骰子干了什么，不可能不可能不可能不可能*N....（逐渐失去高光）
`
var guguTextTable []string
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

	helpWW := `.ww 10a5 // 也可使用.r 10a5替代
.ww 10a5k6m7 // a加骰线 k成功线 m面数
.ww 10 // 骰10a10(默认情况下)
.ww set k6 // 修改成功线为6(当前群)
.ww set a8k6m9 // 修改其他默认设定
.ww set clr // 取消修改`
	cmdWW := CmdItemInfo{
		Name:      "ww",
		ShortHelp: helpWW,
		Help:      "骰池(WOD/无限规则骰点):\n" + helpWW,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			groupAttrs := lo.Must(ctx.Dice.AttrsManager.LoadById(ctx.Group.GroupID))
			switch cmdArgs.GetArgN(1) {
			case "help":
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			case "set":
				arg2 := cmdArgs.GetArgN(2)
				if arg2 == "clr" || arg2 == "clear" {
					groupAttrs.Delete("wodThreshold")
					groupAttrs.Delete("wodPoints")
					groupAttrs.Delete("wodAdd")
					ReplyToSender(ctx, msg, "骰池设定已恢复默认")
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				var texts []string

				reK := regexp.MustCompile(`[kK](\d+)`)
				if m := reK.FindStringSubmatch(arg2); len(m) > 0 {
					if v, err := strconv.ParseInt(m[1], 10, 64); err == nil {
						if v >= 1 {
							groupAttrs.Store("wodThreshold", ds.NewIntVal(ds.IntType(v)))
							texts = append(texts, fmt.Sprintf("成功线k: 已修改为%d", v))
						} else {
							texts = append(texts, "成功线k: 需要至少为1")
						}
					}
				}
				reM := regexp.MustCompile(`[mM](\d+)`)
				if m := reM.FindStringSubmatch(arg2); len(m) > 0 {
					if v, err := strconv.ParseInt(m[1], 10, 64); err == nil {
						if v >= 1 && v <= 2000 {
							groupAttrs.Store("wodPoints", ds.NewIntVal(ds.IntType(v)))
							texts = append(texts, fmt.Sprintf("骰子面数m: 已修改为%d", v))
						} else {
							texts = append(texts, "骰子面数m: 需要在1-2000之间")
						}
					}
				}
				reA := regexp.MustCompile(`[aA](\d+)`)
				if m := reA.FindStringSubmatch(arg2); len(m) > 0 {
					if v, err := strconv.ParseInt(m[1], 10, 64); err == nil {
						if v >= 2 {
							groupAttrs.Store("wodAdd", ds.NewIntVal(ds.IntType(v)))
							texts = append(texts, fmt.Sprintf("加骰线a: 已修改为%d", v))
						} else {
							texts = append(texts, "加骰线a: 需要至少为2")
						}
					}
				}

				if len(texts) == 0 {
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}
				ReplyToSender(ctx, msg, strings.Join(texts, "\n"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			addNum := int64(10)
			if adding, exists := groupAttrs.LoadX("wodAdd"); exists {
				addNumX, _ := adding.ReadInt()
				addNum = int64(addNumX)
			}

			txt := readNumber(cmdArgs.CleanArgs, fmt.Sprintf("a%d", addNum))
			if txt == "" {
				txt = fmt.Sprintf("10a%d", addNum)
				cmdArgs.Args = []string{txt}
			}
			cmdArgs.CleanArgs = txt

			roll := ctx.Dice.CmdMap["roll"]
			ctx.diceExprOverwrite = "10a10"
			return roll.Solve(ctx, msg, cmdArgs)
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
			"rsr":     &cmdRsr,
			"ek":      &cmdEk,
			"ekgen":   &cmdEkgen,
			"dx":      &cmdDX,
			"w":       &cmdWW,
			"ww":      &cmdWW,
			"dxh":     &cmdDX,
			"wh":      &cmdWW,
			"wwh":     &cmdWW,
			"jsr":     &cmdJsr,
			"drl":     &cmdDrl,
			"drlh":    &cmdDrl,
		},
	})
}

func fingerprint(b string) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(b))
	return hash.Sum64()
}