package dice

import (
	"fmt"
	"path/filepath"

	"sealdice-core/message"
	"sealdice-core/utils"
)

type HTTPSimpleMessage struct {
	UID         string `json:"uid"`
	Message     string `json:"message"`
	MessageType string `json:"messageType"`
}

type PlatformAdapterHTTP struct {
	Session       *IMSession
	EndPoint      *EndPointInfo
	RecentMessage []HTTPSimpleMessage
}

func (pa *PlatformAdapterHTTP) SendSegmentToGroup(ctx *MsgContext, groupID string, msg []message.IMessageElement, flag string) {
}

func (pa *PlatformAdapterHTTP) SendSegmentToPerson(ctx *MsgContext, userID string, msg []message.IMessageElement, flag string) {
}

func (pa *PlatformAdapterHTTP) GetGroupInfoAsync(_ string) {}

func (pa *PlatformAdapterHTTP) Serve() int {
	return 0
}

func (pa *PlatformAdapterHTTP) DoRelogin() bool {
	return false
}

func (pa *PlatformAdapterHTTP) SetEnable(_ bool) {}

func (pa *PlatformAdapterHTTP) SendToPerson(ctx *MsgContext, uid string, text string, flag string) {
	sp := utils.SplitLongText(text, 300, utils.DefaultSplitPaginationHint)
	for _, sub := range sp {
		pa.RecentMessage = append(pa.RecentMessage, HTTPSimpleMessage{uid, sub, "private"})
	}
	pa.Session.OnMessageSend(ctx, &Message{
		MessageType: "private",
		Platform:    "UI",
		Message:     text,
		Sender: SenderBase{
			UserID:   pa.EndPoint.UserID,
			Nickname: pa.EndPoint.Nickname,
		},
	}, flag)
}

func (pa *PlatformAdapterHTTP) SendToGroup(ctx *MsgContext, uid string, text string, flag string) {
	sp := utils.SplitLongText(text, 300, utils.DefaultSplitPaginationHint)
	for _, sub := range sp {
		pa.RecentMessage = append(pa.RecentMessage, HTTPSimpleMessage{uid, sub, "group"})
	}
	pa.Session.OnMessageSend(ctx, &Message{
		MessageType: "group",
		Platform:    "UI",
		Message:     text,
		GroupID:     "UI-Group:2001",
		Sender: SenderBase{
			UserID:   pa.EndPoint.UserID,
			Nickname: pa.EndPoint.Nickname,
		},
	}, flag)
}

func (pa *PlatformAdapterHTTP) SendFileToPerson(ctx *MsgContext, uid string, path string, flag string) {
	pa.SendToPerson(ctx, uid, fmt.Sprintf("[尝试发送文件: %s，但不支持]", filepath.Base(path)), flag)
}

func (pa *PlatformAdapterHTTP) SendFileToGroup(ctx *MsgContext, uid string, path string, flag string) {
	pa.SendToGroup(ctx, uid, fmt.Sprintf("[尝试发送文件: %s，但不支持]", filepath.Base(path)), flag)
}

func (pa *PlatformAdapterHTTP) QuitGroup(_ *MsgContext, _ string) {}

func (pa *PlatformAdapterHTTP) SetGroupCardName(_ *MsgContext, _ string) {}

func (pa *PlatformAdapterHTTP) MemberBan(_ string, _ string, _ int64) {}

func (pa *PlatformAdapterHTTP) MemberUnban(_ string, _ string) {}

func (pa *PlatformAdapterHTTP) MemberWholeBan(_ string, _ bool) {}

func (pa *PlatformAdapterHTTP) MemberKick(_ string, _ string) {}

func (pa *PlatformAdapterHTTP) DeleteFriend(_ *MsgContext, _ string) {}

func (pa *PlatformAdapterHTTP) EditMessage(_ *MsgContext, _ string, _ string) {}

func (pa *PlatformAdapterHTTP) RecallMessage(_ *MsgContext, _ string) {}

func (pa *PlatformAdapterHTTP) SendToGroupNotice(_ *MsgContext, _ string, _ string) {}

func (pa *PlatformAdapterHTTP) SendLike(_ *MsgContext, _ string, _ int) {}

func (pa *PlatformAdapterHTTP) SetGroupAdmin(_ *MsgContext, _ string, _ string, _ bool) {}

func (pa *PlatformAdapterHTTP) SetGroupName(_ *MsgContext, _ string, _ string) {}

func (pa *PlatformAdapterHTTP) SetGroupSpecialTitle(_ *MsgContext, _ string, _ string, _ string) {}

func (pa *PlatformAdapterHTTP) SetSelfLongNick(_ *MsgContext, _ string) {}

func (pa *PlatformAdapterHTTP) SharePeer(_ *MsgContext, _ string, _ string) {}

func (pa *PlatformAdapterHTTP) ShareGroup(_ *MsgContext, _ string) {}
