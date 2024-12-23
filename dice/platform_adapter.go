package dice

import "tempestdice/message"

type PlatformAdapter interface {
	Serve() int
	DoRelogin() bool
	SetEnable(enable bool)
	QuitGroup(ctx *MsgContext, ID string)

	SendToPerson(ctx *MsgContext, userID string, text string, flag string)
	SendToGroup(ctx *MsgContext, groupID string, text string, flag string)
	SetGroupCardName(ctx *MsgContext, name string)

	SendSegmentToGroup(ctx *MsgContext, groupID string, msg []message.IMessageElement, flag string)
	SendSegmentToPerson(ctx *MsgContext, userID string, msg []message.IMessageElement, flag string)

	SendFileToPerson(ctx *MsgContext, userID string, path string, flag string)
	SendFileToGroup(ctx *MsgContext, groupID string, path string, flag string)

	MemberBan(groupID string, userID string, duration int64)
	MemberUnban(groupID string, userID string)
	MemberWholeBan(groupID string, enable bool)
	MemberKick(groupID string, userID string)

	GetGroupInfoAsync(groupID string)

	// DeleteFriend 删除好友，目前只有 QQ 平台下的 gocq 和 walleq 实现有这个方法
	DeleteFriend(ctx *MsgContext, id string)

	EditMessage(ctx *MsgContext, msgID string, message string)
	RecallMessage(ctx *MsgContext, msgID string)
	SendToGroupNotice(ctx *MsgContext, groupID string, content string)
	SendLike(ctx *MsgContext, userID string, times int)
	SetGroupAdmin(ctx *MsgContext, groupID string, userID string, enable bool)
	SetGroupName(ctx *MsgContext, groupID string, groupName string)
	SetGroupSpecialTitle(ctx *MsgContext, groupID string, userID string, specialTitle string)
	SetSelfLongNick(ctx *MsgContext, longNick string)
	SharePeer(ctx *MsgContext, groupID string, userID string) //暂时不能用
	ShareGroup(ctx *MsgContext, groupID string)               //暂时不能用
}

// 实现检查
var (
	_ PlatformAdapter = (*PlatformAdapterGocq)(nil)
	_ PlatformAdapter = (*PlatformAdapterHTTP)(nil)
)
