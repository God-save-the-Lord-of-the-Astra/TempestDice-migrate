package dice

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/time/rate"

	ds "github.com/sealdice/dicescript"
)

var (
	sealCodeRe = regexp.MustCompile(`\[(img|图|文本|text|语音|voice|视频|video):(.+?)]`)
	cqCodeRe   = regexp.MustCompile(`\[CQ:.+?]`)
)

func IsCurGroupBotOnByID(session *IMSession, ep *EndPointInfo, messageType string, groupID string) bool {
	// Pinenutn: 总觉得这里还能优化，但是又想不到怎么优化，可恶，要长脑子了
	a := messageType == "group" && session.ServiceAtNew.Exists(groupID)
	if !a {
		return false
	}
	groupInfo, ok := session.ServiceAtNew.Load(groupID)
	if !ok {
		// Pinenutn: 这里是否要打一下日志呢……
		return false
	}
	_, exists := groupInfo.DiceIDActiveMap.Load(ep.UserID)
	return exists
}

func SetBotOffAtGroup(ctx *MsgContext, groupID string) {
	session := ctx.Session
	groupInfo, ok := session.ServiceAtNew.Load(groupID)
	if ok {
		if groupInfo.DiceIDActiveMap == nil {
			groupInfo.DiceIDActiveMap = new(SyncMap[string, bool])
		}

		// TODO: 进行更好的是否变更的检查
		groupInfo.DiceIDActiveMap.Delete(ctx.EndPoint.UserID)
		if groupInfo.DiceIDActiveMap.Len() == 0 {
			groupInfo.Active = false
		}
		groupInfo.UpdatedAtTime = time.Now().Unix()
	}
}

// SetBotOnAtGroup 在群内开启
func SetBotOnAtGroup(ctx *MsgContext, groupID string) *GroupInfo {
	session := ctx.Session
	group, ok := session.ServiceAtNew.Load(groupID)
	if ok {
		if group.DiceIDActiveMap == nil {
			group.DiceIDActiveMap = new(SyncMap[string, bool])
		}
		if group.DiceIDExistsMap == nil {
			group.DiceIDActiveMap = new(SyncMap[string, bool])
		}
		group.DiceIDActiveMap.Store(ctx.EndPoint.UserID, true)
		group.Active = true
	} else {
		// 设定扩展情况
		sort.Sort(ExtDefaultSettingItemSlice(session.Parent.ExtDefaultSettings))
		var extLst []*ExtInfo
		for _, i := range session.Parent.ExtDefaultSettings {
			if i.ExtItem != nil {
				if i.AutoActive {
					extLst = append(extLst, i.ExtItem)
				}
			}
		}

		session.ServiceAtNew.Store(groupID, &GroupInfo{
			Active:           true,
			ActivatedExtList: extLst,
			Players:          new(SyncMap[string, *GroupPlayerInfo]),
			GroupID:          groupID,
			DiceIDActiveMap:  new(SyncMap[string, bool]),
			DiceIDExistsMap:  new(SyncMap[string, bool]),
			CocRuleIndex:     int(session.Parent.DefaultCocRuleIndex),
			UpdatedAtTime:    time.Now().Unix(),
		})
		// TODO: Pinenutn:总觉得这里不太对，但是又觉得合理,GPT也没说怎么改更好一些，求教
		group, _ = session.ServiceAtNew.Load(groupID)
	}

	if group.DiceIDActiveMap == nil {
		group.DiceIDActiveMap = new(SyncMap[string, bool])
	}
	if group.DiceIDExistsMap == nil {
		group.DiceIDExistsMap = new(SyncMap[string, bool])
	}
	if group.BotList == nil {
		group.BotList = new(SyncMap[string, bool])
	}

	group.DiceIDActiveMap.Store(ctx.EndPoint.UserID, true)
	group.UpdatedAtTime = time.Now().Unix()
	return group
}

// GetPlayerInfoBySender 获取玩家群内信息，没有就创建
func GetPlayerInfoBySender(ctx *MsgContext, msg *Message) (*GroupInfo, *GroupPlayerInfo) {
	session := ctx.Session
	var groupID string
	if msg.MessageType == "group" {
		// 群信息
		groupID = msg.GroupID
	} else {
		// 私聊信息 PrivateGroup
		groupID = "PG-" + msg.Sender.UserID
		SetBotOnAtGroup(ctx, groupID)
	}

	// Pinenutn:ServiceAtNew
	groupInfo, ok := session.ServiceAtNew.Load(groupID)
	if !ok {
		groupInfo = SetBotOnAtGroup(ctx, groupID)
	}
	if msg.GuildID != "" {
		groupInfo.GuildID = msg.GuildID
	}
	if msg.ChannelID != "" {
		groupInfo.ChannelID = msg.ChannelID
	}

	p := groupInfo.PlayerGet(ctx.Dice.DBData, msg.Sender.UserID)
	if p == nil {
		p = &GroupPlayerInfo{
			Name:          msg.Sender.Nickname,
			UserID:        msg.Sender.UserID,
			ValueMapTemp:  &ds.ValueMap{},
			UpdatedAtTime: 0, // 新创建时不赋值，这样不会入库保存，减轻数据库负担
		}
		groupInfo.Players.Store(msg.Sender.UserID, p)
	}
	if p.ValueMapTemp == nil {
		p.ValueMapTemp = &ds.ValueMap{}
	}
	p.InGroup = true
	return groupInfo, p
}

func ReplyToSenderRaw(ctx *MsgContext, msg *Message, text string, flag string) {
	inGroup := msg.MessageType == "group"
	if inGroup {
		ReplyGroupRaw(ctx, msg, text, flag)
	} else {
		ReplyPersonRaw(ctx, msg, text, flag)
	}
}

func replyToSenderRawNoCheck(ctx *MsgContext, msg *Message, text string, flag string) {
	inGroup := msg.MessageType == "group"
	if inGroup {
		replyGroupRawNoCheck(ctx, msg, text, flag)
	} else {
		replyPersonRawNoCheck(ctx, msg, text, flag)
	}
}

func ReplyToSender(ctx *MsgContext, msg *Message, text string) {
	ReplyToSenderStartTime := time.Now().UnixMicro()
	go ReplyToSenderRaw(ctx, msg, text, "")
	ReplyToSenderTime := time.Now().UnixMicro() - ReplyToSenderStartTime
	ctx.Dice.Logger.Infof("ReplyToSender 耗时：%dus", ReplyToSenderTime)
}

func ReplyToSenderNoCheck(ctx *MsgContext, msg *Message, text string) {
	go replyToSenderRawNoCheck(ctx, msg, text, "")
}

func ReplyGroupRaw(ctx *MsgContext, msg *Message, text string, flag string) {
	if ctx.AliasPrefixText != "" {
		text = ctx.AliasPrefixText + text
		ctx.AliasPrefixText = ""
	}
	if ctx.DelegateText != "" {
		text = ctx.DelegateText + text
		ctx.DelegateText = ""
	}

	if ctx.Dice.RateLimitEnabled && msg.Platform == "QQ" {
		if !spamCheckPerson(ctx, msg) {
			spamCheckGroup(ctx, msg)
		}
	}

	d := ctx.Dice
	if d != nil {
		d.Logger.Infof("发给(群%s): %s", msg.GroupID, text)
		// 敏感词拦截：回复（群）
		if d.EnableCensor && d.CensorMode == OnlyOutputReply {
			// 先拿掉海豹码和CQ码再检查敏感词
			checkText := sealCodeRe.ReplaceAllString(text, "")
			checkText = cqCodeRe.ReplaceAllString(checkText, "")

			hit, words, needToTerminate, _ := d.CensorMsg(ctx, msg, checkText, text)
			if needToTerminate {
				return
			}
			if hit {
				d.Logger.Infof(
					"拒绝回复命中敏感词「%s」的内容「%s」，原消息「%s」- 来自群(%s)内<%s>(%s)",
					strings.Join(words, "|"),
					text, msg.Message,
					msg.GroupID,
					msg.Sender.Nickname,
					msg.Sender.UserID,
				)
				text = DiceFormatTmpl(ctx, "核心:拦截_完全拦截_发出的消息")
			}
		}
	}
	replyGroupRawNoCheck(ctx, msg, text, flag)
}

func replyGroupRawNoCheck(ctx *MsgContext, msg *Message, text string, flag string) {
	if ctx.AliasPrefixText != "" {
		text = ctx.AliasPrefixText + text
		ctx.AliasPrefixText = ""
	}
	if ctx.DelegateText != "" {
		text = ctx.DelegateText + text
		ctx.DelegateText = ""
	}
	if lenWithoutBase64(text) > 15000 {
		text = "要发送的文本过长"
	}
	if ctx.Group != nil {
		now := time.Now().Unix()
		ctx.Group.RecentDiceSendTime = now
		ctx.Group.UpdatedAtTime = now
	}
	text = strings.TrimSpace(text)
	for _, i := range ctx.SplitText(text) {
		if ctx.EndPoint != nil && ctx.EndPoint.Platform == "QQ" {
			doSleepQQ(ctx)
		}
		ctx.EndPoint.Adapter.SendToGroup(ctx, msg.GroupID, strings.TrimSpace(i), flag)
	}
}

func ReplyGroup(ctx *MsgContext, msg *Message, text string) {
	ReplyGroupStartTime := time.Now().UnixMicro()
	ReplyGroupRaw(ctx, msg, text, "")
	ReplyGroupTime := time.Now().UnixMicro() - ReplyGroupStartTime
	ctx.Dice.Logger.Infof("ReplyGroup 耗时：%dus", ReplyGroupTime)
}

func ReplyPersonRaw(ctx *MsgContext, msg *Message, text string, flag string) {
	if ctx.AliasPrefixText != "" {
		text = ctx.AliasPrefixText + text
		ctx.AliasPrefixText = ""
	}
	if ctx.DelegateText != "" {
		text = ctx.DelegateText + text
		ctx.DelegateText = ""
	}

	if ctx.Dice.RateLimitEnabled && msg.Platform == "QQ" {
		spamCheckPerson(ctx, msg)
	}

	d := ctx.Dice
	if d != nil {
		d.Logger.Infof("发给(帐号%s): %s", msg.Sender.UserID, text)
		// 敏感词拦截：回复（个人）
		if d.EnableCensor && d.CensorMode == OnlyOutputReply {
			// 先拿掉海豹码和CQ码再检查敏感词
			checkText := sealCodeRe.ReplaceAllString(text, "")
			checkText = cqCodeRe.ReplaceAllString(checkText, "")

			hit, words, needToTerminate, _ := d.CensorMsg(ctx, msg, checkText, text)
			if needToTerminate {
				return
			}
			if hit {
				d.Logger.Infof("拒绝回复命中敏感词「%s」的内容「%s」，原消息「%s」- 来自<%s>(%s)",
					strings.Join(words, "|"),
					text,
					msg.Message,
					msg.Sender.Nickname,
					msg.Sender.UserID,
				)
				text = DiceFormatTmpl(ctx, "核心:拦截_完全拦截_发出的消息")
			}
		}
	}
	replyPersonRawNoCheck(ctx, msg, text, flag)
}

func replyPersonRawNoCheck(ctx *MsgContext, msg *Message, text string, flag string) {
	if ctx.AliasPrefixText != "" {
		text = ctx.AliasPrefixText + text
		ctx.AliasPrefixText = ""
	}
	if ctx.DelegateText != "" {
		text = ctx.DelegateText + text
		ctx.DelegateText = ""
	}
	if lenWithoutBase64(text) > 15000 {
		text = "要发送的文本过长"
	}
	text = strings.TrimSpace(text)
	for _, i := range ctx.SplitText(text) {
		if ctx.EndPoint != nil && ctx.EndPoint.Platform == "QQ" {
			doSleepQQ(ctx)
		}
		ctx.EndPoint.Adapter.SendToPerson(ctx, msg.Sender.UserID, strings.TrimSpace(i), flag)
	}
}

// CrossMsgBySearch
// 在 se 中找到第一个平台等于 p 且启用的 EndPointInfo, 并向目标 t 发送消息,
// pr 判断是否为私聊消息
func CrossMsgBySearch(se *IMSession, p, t, txt string, pr bool) bool {
	ep := se.GetEpByPlatform(p)
	if ep == nil {
		return false
	}
	mctx := &MsgContext{
		EndPoint: ep,
		Session:  ep.Session,
		Dice:     ep.Session.Parent,
	}

	if groupInfo, ok := mctx.Session.ServiceAtNew.Load(t); ok {
		mctx.IsCurGroupBotOn = groupInfo.Active
		mctx.Group = groupInfo
	}

	if !pr {
		mctx.MessageType = "group"
		ReplyGroup(mctx, &Message{GroupID: t}, txt)
	} else {
		mctx.IsPrivate = true
		mctx.MessageType = "private"
		ReplyPerson(mctx, &Message{Sender: SenderBase{UserID: t}}, txt)
	}

	return true
}

// TODO: CrossMsgById 用指定 Id 的 EndPoint 发送跨平台消息，现在似乎没有这个需求

func ReplyPerson(ctx *MsgContext, msg *Message, text string) {
	ReplyPersonStartTime := time.Now().UnixMicro()
	ReplyPersonRaw(ctx, msg, text, "")
	ReplyPersonTime := time.Now().UnixMicro() - ReplyPersonStartTime
	ctx.Dice.Logger.Infof("ReplyPerson 耗时：%dus", ReplyPersonTime)
}

func SendFileToSenderRaw(ctx *MsgContext, msg *Message, path string, flag string) {
	inGroup := msg.MessageType == "group"
	if inGroup {
		SendFileToGroupRaw(ctx, msg, path, flag)
	} else {
		SendFileToPersonRaw(ctx, msg, path, flag)
	}
}

func SendFileToPersonRaw(ctx *MsgContext, msg *Message, path string, flag string) {
	if ctx.Dice != nil {
		ctx.Dice.Logger.Infof("发文件给(账号%s): %s", msg.Sender.UserID, path)
	}
	ctx.EndPoint.Adapter.SendFileToPerson(ctx, msg.Sender.UserID, path, flag)
}

func SendFileToGroupRaw(ctx *MsgContext, msg *Message, path string, flag string) {
	if ctx.Dice != nil {
		ctx.Dice.Logger.Infof("发文件给(群%s): %s", msg.GroupID, path)
	}
	ctx.EndPoint.Adapter.SendFileToGroup(ctx, msg.GroupID, path, flag)
}

func MemberBan(ctx *MsgContext, groupID string, userID string, duration int64) {
	MemberBanStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.MemberBan(groupID, userID, duration)
	MemberBanTime := time.Now().UnixMicro() - MemberBanStartTime
	ctx.Dice.Logger.Infof("MemberBan 耗时：%dus", MemberBanTime)
}

func MemberUnban(ctx *MsgContext, groupID string, userID string) {
	MemberUnbanStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.MemberUnban(groupID, userID)
	MemberUnbanTime := time.Now().UnixMicro() - MemberUnbanStartTime
	ctx.Dice.Logger.Infof("MemberBan 耗时：%dus", MemberUnbanTime)
}

func MemberWholeBan(ctx *MsgContext, groupID string, enable bool) {
	MemberWholeBanStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.MemberWholeBan(groupID, enable)
	MemberWholeBanTime := time.Now().UnixMicro() - MemberWholeBanStartTime
	ctx.Dice.Logger.Infof("MemberWholeBan 耗时：%dus", MemberWholeBanTime)
}

func MemberKick(ctx *MsgContext, groupID string, userID string) {
	MemberKickStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.MemberKick(groupID, userID)
	MemberKickTime := time.Now().UnixMicro() - MemberKickStartTime
	ctx.Dice.Logger.Infof("MemberKick 耗时：%dus", MemberKickTime)
}

func GetGroupInfoAsync(ctx *MsgContext, groupID string) {
	GetGroupInfoStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.GetGroupInfoAsync(groupID)
	GetGroupInfoTime := time.Now().UnixMicro() - GetGroupInfoStartTime
	ctx.Dice.Logger.Infof("GetGroupInfo 耗时：%dus", GetGroupInfoTime)
}

func DeleteFriend(ctx *MsgContext, id string) {
	DeleteFriendStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.DeleteFriend(ctx, id)
	DeleteFriendTime := time.Now().UnixMicro() - DeleteFriendStartTime
	ctx.Dice.Logger.Infof("DeleteFriend 耗时：%dus", DeleteFriendTime)
}

func EditMessage(ctx *MsgContext, msgID string, message string) {
	EditMessageStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.EditMessage(ctx, msgID, message)
	EditMessageTime := time.Now().UnixMicro() - EditMessageStartTime
	ctx.Dice.Logger.Infof("EditMessage 耗时：%dus", EditMessageTime)
}

func RecallMessage(ctx *MsgContext, msgID string) {
	RecallMessageStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.RecallMessage(ctx, msgID)
	RecallMessageTime := time.Now().UnixMicro() - RecallMessageStartTime
	ctx.Dice.Logger.Infof("RecallMessage 耗时：%dus", RecallMessageTime)
}

func SendToGroupNotice(ctx *MsgContext, groupID string, content string) {
	SendToGroupNoticeStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SendToGroupNotice(ctx, groupID, content)
	SendToGroupNoticeTime := time.Now().UnixMicro() - SendToGroupNoticeStartTime
	ctx.Dice.Logger.Infof("SendToGroupNotice 耗时：%dus", SendToGroupNoticeTime)
}

func SendLike(ctx *MsgContext, UserID string, times int) {
	SendLikeStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SendLike(ctx, UserID, times)
	SendLikeTime := time.Now().UnixMicro() - SendLikeStartTime
	ctx.Dice.Logger.Infof("SendLike 耗时：%dus", SendLikeTime)
}

func SetGroupAdmin(ctx *MsgContext, groupID string, userID string, enable bool) {
	SetGroupAdminStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SetGroupAdmin(ctx, groupID, userID, enable)
	SetGroupAdminTime := time.Now().UnixMicro() - SetGroupAdminStartTime
	ctx.Dice.Logger.Infof("SetGroupAdmin 耗时：%dus", SetGroupAdminTime)
}

func SetGroupName(ctx *MsgContext, groupID string, name string) {
	SetGroupNameStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SetGroupName(ctx, groupID, name)
	SetGroupNameTime := time.Now().UnixMicro() - SetGroupNameStartTime
	ctx.Dice.Logger.Infof("SetGroupName 耗时：%dus", SetGroupNameTime)
}

func SetGroupSpecialTitle(ctx *MsgContext, groupID string, userID string, specialTitle string) {
	SetGroupSpecialTitleStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SetGroupSpecialTitle(ctx, groupID, userID, specialTitle)
	SetGroupSpecialTitleTime := time.Now().UnixMicro() - SetGroupSpecialTitleStartTime
	ctx.Dice.Logger.Infof("SetGroupSpecialTitle 耗时：%dus", SetGroupSpecialTitleTime)
}

func SetSelfLongNick(ctx *MsgContext, longNick string) {
	SetSelfLongNickStartTime := time.Now().UnixMicro()
	ctx.EndPoint.Adapter.SetSelfLongNick(ctx, longNick)
	SetSelfLongNickTime := time.Now().UnixMicro() - SetSelfLongNickStartTime
	ctx.Dice.Logger.Infof("SetSelfLongNick 耗时：%dus", SetSelfLongNickTime)
}

func SharePeer(ctx *MsgContext, groupID string, userID string) {
	ctx.EndPoint.Adapter.SharePeer(ctx, groupID, userID)
}

func ShareGroup(ctx *MsgContext, groupID string) {
	ctx.EndPoint.Adapter.ShareGroup(ctx, groupID)
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}

func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func CompatibleReplace(ctx *MsgContext, s string) string {
	s = ctx.TranslateSplit(s)

	// 匹配 #{DRAW-$1}, 其中$1执行最短匹配且允许左侧右侧各有一个花括号
	// #{DRAW-aaa} => aaa
	// #{DRAW-{aaa} => {aaa
	// #{DRAW-aaa}} => aaa}
	// #{DRAW-{aaa}} => {aaa}
	// 这允许在牌组名中使用不含空格的表达式(主要是为了变量)
	re := regexp.MustCompile(`#\{DRAW-(\{?\S+?\}?)\}`)
	s = re.ReplaceAllString(s, "###DRAW-$1###")

	if ctx != nil {
		s = DeckRewrite(s, func(deckName string) string {
			// 如果牌组名中含有表达式, 在此进行求值
			// 不含表达式也无妨, 求值完还是原来的字符串
			r, _, err := DiceExprTextBase(ctx, deckName, RollExtraFlags{})
			if err == nil {
				deckName = r.ToString()
			}

			exists, result, err := deckDraw(ctx, deckName, false)
			if !exists {
				return "<%未知牌组-" + deckName + "%>"
			}
			if err != nil {
				return "<%抽取错误-" + deckName + "%>"
			}
			return result
		})
	}
	return s
}

func FormatDiceID(ctx *MsgContext, id interface{}, isGroup bool) string {
	prefix := ctx.EndPoint.Platform
	if isGroup {
		prefix += "-Group"
	}
	return fmt.Sprintf("%s:%v", prefix, id)
}

func spamCheckPerson(ctx *MsgContext, msg *Message) bool {
	if ctx.SpamCheckedPerson {
		return false
	}

	// 同一个 ctx 只需检查一次
	defer func() {
		ctx.SpamCheckedPerson = true
	}()

	if ctx.PrivilegeLevel >= 100 {
		return false
	}

	if ctx.Player.RateLimiter == nil {
		ctx.Player.RateLimitWarned = false
		if ctx.Dice.PersonalReplenishRateStr == "" {
			ctx.Dice.PersonalReplenishRateStr = "@every 3s"
			ctx.Dice.PersonalReplenishRate = rate.Every(time.Second * 3)
		}
		if ctx.Dice.PersonalBurst == 0 {
			ctx.Dice.PersonalBurst = 3
		}
		ctx.Player.RateLimiter = rate.NewLimiter(
			ctx.Dice.PersonalReplenishRate,
			int(ctx.Dice.PersonalBurst),
		)
	}

	if ctx.Player.RateLimiter.Allow() {
		ctx.Player.RateLimitWarned = false
		return false
	}

	if ctx.Player.RateLimitWarned {
		ctx.Dice.BanList.AddScoreByCommandSpam(ctx.Player.UserID, msg.GroupID, ctx)
	} else {
		ctx.Player.RateLimitWarned = true
		replyToSenderRawNoCheck(
			ctx, msg,
			DiceFormatTmpl(ctx, "核心:刷屏_警告内容_个人"),
			"",
		)
	}

	return true
}

func spamCheckGroup(ctx *MsgContext, msg *Message) bool {
	if ctx.SpamCheckedGroup {
		return false
	}

	// 同一个 ctx 只需检查一次
	defer func() {
		ctx.SpamCheckedGroup = true
	}()

	// Skip privileged groups
	for _, g := range ctx.Dice.DiceMasters {
		if ctx.Group.GroupID == g {
			return false
		}
	}

	if ctx.Group.RateLimiter == nil {
		ctx.Group.RateLimitWarned = false
		if ctx.Dice.GroupReplenishRateStr == "" {
			ctx.Dice.GroupReplenishRateStr = "@every 3s"
			ctx.Dice.GroupReplenishRate = rate.Every(time.Second * 3)
		}
		if ctx.Dice.GroupBurst == 0 {
			ctx.Dice.GroupBurst = 3
		}
		ctx.Group.RateLimiter = rate.NewLimiter(
			ctx.Dice.GroupReplenishRate,
			int(ctx.Dice.GroupBurst),
		)
	}

	if ctx.Group.RateLimiter.Allow() {
		ctx.Group.RateLimitWarned = false
		return false
	}

	// If not allow
	if ctx.Group.RateLimitWarned {
		ctx.Dice.BanList.AddScoreByCommandSpam(ctx.Group.GroupID, msg.GroupID, ctx)
	} else {
		ctx.Group.RateLimitWarned = true
		replyToSenderRawNoCheck(
			ctx, msg,
			DiceFormatTmpl(ctx, "核心:刷屏_警告内容_群组"),
			"",
		)
	}

	return true
}

func lenWithoutBase64(text string) int {
	re := regexp.MustCompile(`base64://[A-Za-z0-9+/=]+`)
	croppedText := re.ReplaceAllString(text, "")
	return len(croppedText)
}
