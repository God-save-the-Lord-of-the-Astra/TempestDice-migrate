package dice

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func LuaVarInitWithoutArgs(LuaVM *lua.LState, d *Dice, ctx *MsgContext, msg *Message) {
	msgTable := LuaVM.NewTable()

	// 设置Message的字段
	msgTable.RawSetString("Time", lua.LNumber(msg.Time))
	msgTable.RawSetString("MessageType", lua.LString(msg.MessageType))
	msgTable.RawSetString("GroupID", lua.LString(msg.GroupID))
	msgTable.RawSetString("GuildID", lua.LString(msg.GuildID))
	msgTable.RawSetString("ChannelID", lua.LString(msg.ChannelID))
	msgTable.RawSetString("Message", lua.LString(msg.Message))
	msgTable.RawSetString("RawID", lua.LString(fmt.Sprintf("%v", msg.RawID)))
	msgTable.RawSetString("Platform", lua.LString(msg.Platform))
	msgTable.RawSetString("GroupName", lua.LString(msg.GroupName))
	msgTable.RawSetString("TmpUID", lua.LString(msg.TmpUID))

	// 设置Sender的字段
	senderTable := LuaVM.NewTable()
	senderTable.RawSetString("Nickname", lua.LString(msg.Sender.Nickname))
	senderTable.RawSetString("UserID", lua.LString(msg.Sender.UserID))

	// 将senderTable添加到msgTable中
	msgTable.RawSetString("sender", senderTable)

	//----------------------------------------------------------------
	ctxTable := LuaVM.NewTable()
	LuaVM.SetField(ctxTable, "MessageType", lua.LString(ctx.MessageType))
	LuaVM.SetField(ctxTable, "IsCurGroupBotOn", lua.LBool(ctx.IsCurGroupBotOn))
	LuaVM.SetField(ctxTable, "IsPrivate", lua.LBool(ctx.IsPrivate))
	LuaVM.SetField(ctxTable, "CommandID", lua.LNumber(ctx.CommandID))
	LuaVM.SetField(ctxTable, "PrivilegeLevel", lua.LNumber(ctx.PrivilegeLevel))
	LuaVM.SetField(ctxTable, "GroupRoleLevel", lua.LNumber(ctx.GroupRoleLevel))
	LuaVM.SetField(ctxTable, "DelegateText", lua.LString(ctx.DelegateText))
	LuaVM.SetField(ctxTable, "AliasPrefixText", lua.LString(ctx.AliasPrefixText))

	// Group info as a nested table
	if ctx.Group != nil {
		groupTable := LuaVM.NewTable()
		LuaVM.SetField(groupTable, "GroupID", lua.LString(ctx.Group.GroupID))
		LuaVM.SetField(groupTable, "GuildID", lua.LString(ctx.Group.GuildID))
		LuaVM.SetField(groupTable, "ChannelID", lua.LString(ctx.Group.ChannelID))
		LuaVM.SetField(groupTable, "GroupName", lua.LString(ctx.Group.GroupName))
		LuaVM.SetField(groupTable, "RecentDiceSendTime", lua.LNumber(ctx.Group.RecentDiceSendTime))
		LuaVM.SetField(groupTable, "ShowGroupWelcome", lua.LBool(ctx.Group.ShowGroupWelcome))
		LuaVM.SetField(groupTable, "GroupWelcomeMessage", lua.LString(ctx.Group.GroupWelcomeMessage))
		LuaVM.SetField(groupTable, "EnteredTime", lua.LNumber(ctx.Group.EnteredTime))
		LuaVM.SetField(groupTable, "InviteUserID", lua.LString(ctx.Group.InviteUserID))
		LuaVM.SetField(groupTable, "TmpPlayerNum", lua.LNumber(ctx.Group.TmpPlayerNum))
		LuaVM.SetField(groupTable, "UpdatedAtTime", lua.LNumber(ctx.Group.UpdatedAtTime))
		LuaVM.SetField(groupTable, "DefaultHelpGroup", lua.LString(ctx.Group.DefaultHelpGroup))
		LuaVM.SetField(ctxTable, "Group", groupTable)
	}

	// Player info as a nested table
	if ctx.Player != nil {
		playerTable := LuaVM.NewTable()
		LuaVM.SetField(playerTable, "Name", lua.LString(ctx.Player.Name))
		LuaVM.SetField(playerTable, "UserID", lua.LString(ctx.Player.UserID))
		LuaVM.SetField(playerTable, "InGroup", lua.LBool(ctx.Player.InGroup))
		LuaVM.SetField(playerTable, "LastCommandTime", lua.LNumber(ctx.Player.LastCommandTime))
		LuaVM.SetField(playerTable, "RateLimitWarned", lua.LBool(ctx.Player.RateLimitWarned))
		LuaVM.SetField(playerTable, "AutoSetNameTemplate", lua.LString(ctx.Player.AutoSetNameTemplate))
		LuaVM.SetField(playerTable, "DiceSideNum", lua.LNumber(ctx.Player.DiceSideNum))
		LuaVM.SetField(playerTable, "UpdatedAtTime", lua.LNumber(ctx.Player.UpdatedAtTime))
		LuaVM.SetField(playerTable, "RecentUsedTime", lua.LNumber(ctx.Player.RecentUsedTime))
		LuaVM.SetField(ctxTable, "Player", playerTable)
	}

	if ctx.EndPoint != nil {
		endPointTable := LuaVM.NewTable()
		LuaVM.SetField(endPointTable, "Name", lua.LString(ctx.EndPoint.ID))
		LuaVM.SetField(endPointTable, "Nickname", lua.LString(ctx.EndPoint.Nickname))
		LuaVM.SetField(endPointTable, "UserID", lua.LString(ctx.EndPoint.UserID))
		LuaVM.SetField(endPointTable, "GroupNum", lua.LNumber(ctx.EndPoint.GroupNum))
		LuaVM.SetField(endPointTable, "State", lua.LString(ctx.EndPoint.State))
		LuaVM.SetField(endPointTable, "CmdExecutedNum", lua.LNumber(ctx.EndPoint.CmdExecutedNum))
		LuaVM.SetField(endPointTable, "CmdExecutedLastTime", lua.LNumber(ctx.EndPoint.CmdExecutedLastTime))
		LuaVM.SetField(endPointTable, "OnlineTotalTime", lua.LNumber(ctx.EndPoint.OnlineTotalTime))
		LuaVM.SetField(endPointTable, "Platform", lua.LString(ctx.EndPoint.Platform))
		LuaVM.SetField(endPointTable, "RelWorkDir", lua.LString(ctx.EndPoint.RelWorkDir))
		LuaVM.SetField(endPointTable, "Enable", lua.LBool(ctx.EndPoint.Enable))
		LuaVM.SetField(endPointTable, "ProtocolType", lua.LString(ctx.EndPoint.ProtocolType))
		LuaVM.SetField(endPointTable, "IsPublic", lua.LBool(ctx.EndPoint.IsPublic))
		LuaVM.SetField(ctxTable, "EndPoint", endPointTable)
	}

	//----------------------------------------------------------------
	msgUD := LuaVM.NewUserData()
	msgUD.Value = msg
	msgMeta := LuaVM.NewTypeMetatable("Message")
	msgUD.Metatable = LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"__index": func(LuaVM *lua.LState) int {
			LuaVM.Push(msgTable)
			return 1
		},
	})
	LuaVM.SetGlobal("Message", msgMeta)
	LuaVM.SetField(msgMeta, "__index", LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"Time": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LNumber(msg.Time))
			return 1
		},
		"MessageType": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.MessageType))
			return 1
		},
		"GroupID": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.GroupID))
			return 1
		},
		"GuildID": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.GuildID))
			return 1
		},
		"ChannelID": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.ChannelID))
			return 1
		},
		"Message": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.Message))
			return 1
		},
		"Platform": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.Platform))
			return 1
		},
		"GroupName": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			LuaVM.Push(lua.LString(msg.GroupName))
			return 1
		},
		"Sender": func(LuaVM *lua.LState) int {
			msg := LuaVM.CheckUserData(1).Value.(*Message)
			senderTable := LuaVM.NewTable()
			senderTable.RawSetString("Nickname", lua.LString(msg.Sender.Nickname))
			senderTable.RawSetString("UserID", lua.LString(msg.Sender.UserID))
			LuaVM.Push(senderTable)
			return 1
		},
	}))

	LuaVM.SetGlobal("msg", msgUD)

	//----------------------------------------------------------------
	ctxUD := LuaVM.NewUserData()
	ctxUD.Value = ctx
	ctxMeta := LuaVM.NewTypeMetatable("MsgContext")
	ctxUD.Metatable = LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"__index": func(LuaVM *lua.LState) int {
			LuaVM.Push(ctxTable)
			return 1
		},
	})
	LuaVM.SetGlobal("MsgContext", ctxMeta)
	LuaVM.SetField(ctxMeta, "__index", LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"MessageType": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LString(ctx.MessageType))
			return 1
		},
		"IsCurGroupBotOn": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LBool(ctx.IsCurGroupBotOn))
			return 1
		},
		"IsPrivate": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LBool(ctx.IsPrivate))
			return 1
		},
		"CommandID": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LNumber(ctx.CommandID))
			return 1
		},
		"PrivilegeLevel": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LNumber(ctx.PrivilegeLevel))
			return 1
		},
		"GroupRoleLevel": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LNumber(ctx.GroupRoleLevel))
			return 1
		},
		"DelegateText": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LString(ctx.DelegateText))
			return 1
		},
		"AliasPrefixText": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			LuaVM.Push(lua.LString(ctx.AliasPrefixText))
			return 1
		},
		"Group": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			if ctx.Group != nil {
				groupTable := LuaVM.NewTable()
				LuaVM.SetField(groupTable, "GroupID", lua.LString(ctx.Group.GroupID))
				LuaVM.SetField(groupTable, "GuildID", lua.LString(ctx.Group.GuildID))
				LuaVM.SetField(groupTable, "ChannelID", lua.LString(ctx.Group.ChannelID))
				LuaVM.SetField(groupTable, "GroupName", lua.LString(ctx.Group.GroupName))
				LuaVM.SetField(groupTable, "RecentDiceSendTime", lua.LNumber(ctx.Group.RecentDiceSendTime))
				LuaVM.SetField(groupTable, "ShowGroupWelcome", lua.LBool(ctx.Group.ShowGroupWelcome))
				LuaVM.SetField(groupTable, "GroupWelcomeMessage", lua.LString(ctx.Group.GroupWelcomeMessage))
				LuaVM.SetField(groupTable, "EnteredTime", lua.LNumber(ctx.Group.EnteredTime))
				LuaVM.SetField(groupTable, "InviteUserID", lua.LString(ctx.Group.InviteUserID))
				LuaVM.Push(groupTable)
			} else {
				LuaVM.Push(lua.LNil)
			}
			return 1
		},
		"Player": func(LuaVM *lua.LState) int {
			ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
			if ctx.Player != nil {
				playerTable := LuaVM.NewTable()
				LuaVM.SetField(playerTable, "Name", lua.LString(ctx.Player.Name))
				LuaVM.SetField(playerTable, "UserID", lua.LString(ctx.Player.UserID))
				LuaVM.SetField(playerTable, "LastCommandTime", lua.LNumber(ctx.Player.LastCommandTime))
				LuaVM.SetField(playerTable, "AutoSetNameTemplate", lua.LString(ctx.Player.AutoSetNameTemplate))
				LuaVM.Push(playerTable)
			} else {
				LuaVM.Push(lua.LNil)
			}
			return 1
		},
	}))

	LuaVM.SetGlobal("ctx", ctxUD)

	//----------------------------------------------------------------
	// Shiki散装变量兼容
	ShikiMsgTable := LuaVM.NewTable()
	ShikiMsgFromQQ := strings.ReplaceAll(msg.Sender.UserID, "QQ:", "")
	ShikiMsgFromGroup := strings.ReplaceAll(msg.GroupID, "QQ-Group:", "")
	ShikiMsgFromUID, _ := strconv.Atoi(strings.ReplaceAll(msg.Sender.UserID, "QQ:", ""))
	ShikiMsgFromGID, _ := strconv.Atoi(strings.ReplaceAll(msg.GroupID, "QQ-Group:", ""))
	ShikiMsgTimeStamp := msg.Time
	ShikiMsgID := fmt.Sprintf("%v", msg.RawID)

	// Register Shiki variables
	LuaVM.SetField(ShikiMsgTable, "fromQQ", lua.LString(ShikiMsgFromQQ))
	LuaVM.SetField(ShikiMsgTable, "fromGroup", lua.LString(ShikiMsgFromGroup))
	LuaVM.SetField(ShikiMsgTable, "fromUID", lua.LNumber(ShikiMsgFromUID))
	LuaVM.SetField(ShikiMsgTable, "fromGID", lua.LNumber(ShikiMsgFromGID))
	LuaVM.SetField(ShikiMsgTable, "timestamp", lua.LNumber(ShikiMsgTimeStamp))
	LuaVM.SetField(ShikiMsgTable, "msgid", lua.LString(ShikiMsgID))

	LuaVM.SetGlobal("shikimsg", ShikiMsgTable)
	//----------------------------------------------------------------
	// Dream 散装变量兼容
	DreamMsgGroup_ID := strings.ReplaceAll(msg.GroupID, "QQ-Group:", "")
	DreamMsgSender_ID := strings.ReplaceAll(msg.Sender.UserID, "QQ:", "")
	DreamMsgIsGroup := !ctx.IsPrivate
	DreamMsgFromDiceName := ctx.EndPoint.Nickname
	//DreamMsgGroup_Nick := ctx.Group.GroupName
	DreamMsgSender_Nick := ctx.Player.Name
	DreamMsgSender_Jrrp, _ := VarGetValueInt64(ctx, "$t人品")
	DreamMsgTable := LuaVM.NewTable()
	LuaVM.SetField(DreamMsgTable, "fromGroup", lua.LString(DreamMsgGroup_ID))
	LuaVM.SetField(DreamMsgTable, "fromQQ", lua.LString(DreamMsgSender_ID))
	LuaVM.SetField(DreamMsgTable, "isGroup", lua.LBool(DreamMsgIsGroup))
	LuaVM.SetField(DreamMsgTable, "fromDiceName", lua.LString(DreamMsgFromDiceName))
	LuaVM.SetField(DreamMsgTable, "fromNick", lua.LString(DreamMsgSender_Nick))
	LuaVM.SetField(DreamMsgTable, "fromJrrp", lua.LNumber(DreamMsgSender_Jrrp))
	LuaVM.SetGlobal("dreammsg", DreamMsgTable)
}

func LuaVarInit(LuaVM *lua.LState, d *Dice, ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) {
	LuaVarInitWithoutArgs(LuaVM, d, ctx, msg)
	//----------------------------------------------------------------
	cmdArgsTable := LuaVM.NewTable()
	// 注册基本字段
	LuaVM.SetField(cmdArgsTable, "Command", lua.LString(cmdArgs.Command))

	// 注册切片类型字段 (Args)
	argsTable := LuaVM.NewTable()
	for _, arg := range cmdArgs.Args {
		argsTable.Append(lua.LString(arg))
	}
	LuaVM.SetField(cmdArgsTable, "Args", argsTable)

	// 注册结构体数组 (Kwargs)
	kwargsTable := LuaVM.NewTable()
	for _, kwarg := range cmdArgs.Kwargs {
		kwargTable := LuaVM.NewTable()
		LuaVM.SetField(kwargTable, "Name", lua.LString(kwarg.Name))
		LuaVM.SetField(kwargTable, "ValueExists", lua.LBool(kwarg.ValueExists))
		LuaVM.SetField(kwargTable, "Value", lua.LString(kwarg.Value))
		LuaVM.SetField(kwargTable, "AsBool", lua.LBool(kwarg.AsBool))
		kwargsTable.Append(kwargTable)
	}
	LuaVM.SetField(cmdArgsTable, "Kwargs", kwargsTable)

	// 注册结构体数组 (At)
	atInfoTable := LuaVM.NewTable()
	for _, at := range cmdArgs.At {
		atTable := LuaVM.NewTable()
		LuaVM.SetField(atTable, "UserID", lua.LString(at.UserID))
		atInfoTable.Append(atTable)
	}
	LuaVM.SetField(cmdArgsTable, "At", atInfoTable)

	// 注册其他基本字段
	LuaVM.SetField(cmdArgsTable, "RawArgs", lua.LString(cmdArgs.RawArgs))
	LuaVM.SetField(cmdArgsTable, "AmIBeMentioned", lua.LBool(cmdArgs.AmIBeMentioned))
	LuaVM.SetField(cmdArgsTable, "AmIBeMentionedFirst", lua.LBool(cmdArgs.AmIBeMentionedFirst))
	LuaVM.SetField(cmdArgsTable, "SomeoneBeMentionedButNotMe", lua.LBool(cmdArgs.SomeoneBeMentionedButNotMe))
	LuaVM.SetField(cmdArgsTable, "IsSpaceBeforeArgs", lua.LBool(cmdArgs.IsSpaceBeforeArgs))
	LuaVM.SetField(cmdArgsTable, "CleanArgs", lua.LString(cmdArgs.CleanArgs))
	LuaVM.SetField(cmdArgsTable, "SpecialExecuteTimes", lua.LNumber(cmdArgs.SpecialExecuteTimes))
	LuaVM.SetField(cmdArgsTable, "RawText", lua.LString(cmdArgs.RawText))

	//----------------------------------------------------------------
	cmdArgsUD := LuaVM.NewUserData()
	cmdArgsUD.Value = cmdArgs
	cmdArgsMeta := LuaVM.NewTypeMetatable("CmdArgs")
	cmdArgsUD.Metatable = LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"__index": func(LuaVM *lua.LState) int {
			LuaVM.Push(cmdArgsTable)
			return 1
		},
	})
	LuaVM.SetGlobal("CmdArgs", cmdArgsMeta)
	LuaVM.SetField(cmdArgsMeta, "__index", LuaVM.SetFuncs(LuaVM.NewTable(), map[string]lua.LGFunction{
		"Command": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LString(cmdArgs.Command))
			return 1
		},
		"Args": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			argsTable := LuaVM.NewTable()
			for _, arg := range cmdArgs.Args {
				argsTable.Append(lua.LString(arg))
			}
			LuaVM.Push(argsTable)
			return 1
		},
		"Kwargs": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			kwargsTable := LuaVM.NewTable()
			for _, kwarg := range cmdArgs.Kwargs {
				kwargTable := LuaVM.NewTable()
				LuaVM.SetField(kwargTable, "Name", lua.LString(kwarg.Name))
				LuaVM.SetField(kwargTable, "ValueExists", lua.LBool(kwarg.ValueExists))
				LuaVM.SetField(kwargTable, "Value", lua.LString(kwarg.Value))
				LuaVM.SetField(kwargTable, "AsBool", lua.LBool(kwarg.AsBool))
				kwargsTable.Append(kwargTable)
			}
			LuaVM.Push(kwargsTable)
			return 1
		},
		"At": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			atInfoTable := LuaVM.NewTable()
			for _, at := range cmdArgs.At {
				atTable := LuaVM.NewTable()
				LuaVM.SetField(atTable, "UserID", lua.LString(at.UserID))
				atInfoTable.Append(atTable)
			}
			LuaVM.Push(atInfoTable)
			return 1
		},
		"RawArgs": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LString(cmdArgs.RawArgs))
			return 1
		},
		"AmIBeMentioned": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LBool(cmdArgs.AmIBeMentioned))
			return 1
		},
		"AmIBeMentionedFirst": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LBool(cmdArgs.AmIBeMentionedFirst))
			return 1
		},
		"SomeoneBeMentionedButNotMe": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LBool(cmdArgs.SomeoneBeMentionedButNotMe))
			return 1
		},
		"IsSpaceBeforeArgs": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LBool(cmdArgs.IsSpaceBeforeArgs))
			return 1
		},
		"CleanArgs": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LString(cmdArgs.CleanArgs))
			return 1
		},
		"SpecialExecuteTimes": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LNumber(cmdArgs.SpecialExecuteTimes))
			return 1
		},
		"RawText": func(LuaVM *lua.LState) int {
			cmdArgs := LuaVM.CheckUserData(1).Value.(*CmdArgs)
			LuaVM.Push(lua.LString(cmdArgs.RawText))
			return 1
		},
	}))

	LuaVM.SetGlobal("cmdArgs", cmdArgsUD)

	DiceUD := LuaVM.NewUserData()
	DiceUD.Value = d
	DiceMeta := LuaVM.NewTypeMetatable("Dice")
	LuaVM.SetGlobal("Dice", DiceMeta)
	LuaVM.SetGlobal("d", DiceUD)

	//----------------------------------------------------------------
	// Shiki散装变量兼容
	ShikiMsgTable := LuaVM.NewTable()
	ShikiMsgFromQQ := strings.ReplaceAll(msg.Sender.UserID, "QQ:", "")
	ShikiMsgFromGroup := strings.ReplaceAll(msg.GroupID, "QQ-Group:", "")
	ShikiMsgFromUID, _ := strconv.Atoi(strings.ReplaceAll(msg.Sender.UserID, "QQ:", ""))
	ShikiMsgFromGID, _ := strconv.Atoi(strings.ReplaceAll(msg.GroupID, "QQ-Group:", ""))
	ShikiMsgFromMsg := cmdArgs.RawText
	ShikiMsgTimeStamp := msg.Time
	ShikiMsgSuffix := cmdArgs.RawArgs
	ShikiMsgPrefix := strings.TrimSpace(strings.ReplaceAll(cmdArgs.RawText, cmdArgs.RawArgs, ""))
	ShikiMsgID := fmt.Sprintf("%v", msg.RawID)
	ShikiMsgCmdTable := cmdArgs.Args

	// Register Shiki variables
	LuaVM.SetField(ShikiMsgTable, "fromQQ", lua.LString(ShikiMsgFromQQ))
	LuaVM.SetField(ShikiMsgTable, "fromGroup", lua.LString(ShikiMsgFromGroup))
	LuaVM.SetField(ShikiMsgTable, "fromUID", lua.LNumber(ShikiMsgFromUID))
	LuaVM.SetField(ShikiMsgTable, "fromGID", lua.LNumber(ShikiMsgFromGID))
	LuaVM.SetField(ShikiMsgTable, "fromMsg", lua.LString(ShikiMsgFromMsg))
	LuaVM.SetField(ShikiMsgTable, "suffix", lua.LString(ShikiMsgSuffix))
	LuaVM.SetField(ShikiMsgTable, "prefix", lua.LString(ShikiMsgPrefix))
	LuaVM.SetField(ShikiMsgTable, "timestamp", lua.LNumber(ShikiMsgTimeStamp))
	LuaVM.SetField(ShikiMsgTable, "msgid", lua.LString(ShikiMsgID))
	MsgCmdTable := LuaVM.NewTable()
	for _, arg := range ShikiMsgCmdTable {
		MsgCmdTable.Append(lua.LString(arg))
	}
	LuaVM.SetField(ShikiMsgTable, "CmdTab", MsgCmdTable)
	LuaVM.SetGlobal("shikimsg", ShikiMsgTable)
	//----------------------------------------------------------------
	// Dream 散装变量兼容
	DreamMsgGroup_ID := strings.ReplaceAll(msg.GroupID, "QQ-Group:", "")
	DreamMsgSender_ID := strings.ReplaceAll(msg.Sender.UserID, "QQ:", "")
	DreamMsgIsGroup := !ctx.IsPrivate
	DreamMsgIsAtMe := cmdArgs.AmIBeMentionedFirst
	DreamMsgFromDiceName := ctx.EndPoint.Nickname
	DreamMsgfromParas := cmdArgs.RawArgs
	DreamMsgCommandThis := strings.TrimSpace(strings.ReplaceAll(cmdArgs.RawText, cmdArgs.RawArgs, ""))
	//DreamMsgGroup_Nick := ctx.Group.GroupName
	DreamMsgSender_Nick := ctx.Player.Name
	DreamMsgMessage_Text := cmdArgs.RawText
	DreamMsgSender_Jrrp, _ := VarGetValueInt64(ctx, "$t人品")
	DreamMsgTable := LuaVM.NewTable()
	LuaVM.SetField(DreamMsgTable, "fromGroup", lua.LString(DreamMsgGroup_ID))
	LuaVM.SetField(DreamMsgTable, "fromQQ", lua.LString(DreamMsgSender_ID))
	LuaVM.SetField(DreamMsgTable, "isGroup", lua.LBool(DreamMsgIsGroup))
	LuaVM.SetField(DreamMsgTable, "isAtMe", lua.LBool(DreamMsgIsAtMe))
	LuaVM.SetField(DreamMsgTable, "fromDiceName", lua.LString(DreamMsgFromDiceName))
	LuaVM.SetField(DreamMsgTable, "fromParas", lua.LString(DreamMsgfromParas))
	LuaVM.SetField(DreamMsgTable, "commandThis", lua.LString(DreamMsgCommandThis))
	LuaVM.SetField(DreamMsgTable, "fromNick", lua.LString(DreamMsgSender_Nick))
	LuaVM.SetField(DreamMsgTable, "fromJrrp", lua.LNumber(DreamMsgSender_Jrrp))
	LuaVM.SetField(DreamMsgTable, "fromMsg", lua.LString(DreamMsgMessage_Text))
	DreamMsgParaTable := LuaVM.NewTable()
	for _, arg := range cmdArgs.Args {
		DreamMsgParaTable.Append(lua.LString(arg))
	}
	LuaVM.SetField(DreamMsgTable, "ParaTable", DreamMsgParaTable)
	LuaVM.SetGlobal("dreammsg", DreamMsgTable)
}

// ----------------------------------------------------------------
func luaVarSetValueStr(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	v := LuaVM.CheckString(3)
	VarSetValueStr(ctx, s, v)
	return 0 // 返回 0 表示无返回值
}

func luaVarSetValueInt(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	v := LuaVM.CheckInt64(3)
	VarSetValueInt64(ctx, s, v)
	return 0 // 返回 0 表示无返回值
}

func luaVarDelValue(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	VarDelValue(ctx, s)
	return 0 // 返回 0 表示无返回值
}

func luaVarGetValueInt(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	res, exists := VarGetValueInt64(ctx, s)
	if !exists {
		return 0 // 返回 0 表示没有值
	}
	LuaVM.Push(lua.LNumber(res)) // 推送结果到 Lua 栈
	return 1                     // 返回 1 表示成功
}

func luaVarGetValueStr(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	res, exists := VarGetValueStr(ctx, s)
	if !exists {
		return 0 // 返回 0 表示没有值
	}
	LuaVM.Push(lua.LString(res)) // 推送结果到 Lua 栈
	return 1                     // 返回 1 表示成功
}

func luaAddBan(LuaVM *lua.LState) int {
	id := LuaVM.CheckString(1)
	d := LuaVM.CheckUserData(2).Value.(*Dice)
	place := LuaVM.CheckString(3)
	reason := LuaVM.CheckString(4)
	ctx := LuaVM.CheckUserData(5).Value.(*MsgContext)
	d.BanList.AddScoreBase(id, d.BanList.ThresholdBan, place, reason, ctx)
	d.BanList.SaveChanged(d)
	return 1 // 返回 1 表示成功
}

func luaAddTrust(LuaVM *lua.LState) int {
	d := LuaVM.CheckUserData(1).Value.(*Dice)
	id := LuaVM.CheckString(2)
	place := LuaVM.CheckString(3)
	reason := LuaVM.CheckString(4)
	d.BanList.SetTrustByID(id, place, reason)
	d.BanList.SaveChanged(d)
	return 1 // 返回 1 表示成功
}

func luaRemoveBan(LuaVM *lua.LState) int {
	d := LuaVM.CheckUserData(1).Value.(*Dice)
	id := LuaVM.CheckString(2)
	_, ok := d.BanList.GetByID(id)
	if !ok {
		return 0 // 返回 0 表示没有值
	}
	d.BanList.DeleteByID(d, id)
	return 1 // 返回 1 表示成功
}

func luaReplyGroup(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(2).Value.(*Message)
	text := LuaVM.CheckString(3)
	ReplyGroup(ctx, msg, text)
	return 1 // 返回 1 表示成功
}

func luaReplyPerson(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(2).Value.(*Message)
	text := LuaVM.CheckString(3)
	ReplyPerson(ctx, msg, text)
	return 1 // 返回 1 表示成功
}

func luaReplyToSender(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(2).Value.(*Message)
	text := LuaVM.CheckString(3)
	ReplyToSender(ctx, msg, text)
	return 1 // 返回 1 表示成功
}
func luaMemberBan(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	userID := LuaVM.CheckString(3)
	duration := LuaVM.CheckInt64(4)
	MemberBan(ctx, groupID, userID, duration)
	return 1 // 返回 1 表示成功
}
func luaMemberUnban(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	userID := LuaVM.CheckString(3)
	MemberUnban(ctx, groupID, userID)
	return 1
}
func luaMemberWholeBan(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	enable := LuaVM.CheckBool(3)
	MemberWholeBan(ctx, groupID, enable)
	return 1 // 返回 1 表示成功
}
func luaMemberKick(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	userID := LuaVM.CheckString(3)
	MemberKick(ctx, groupID, userID)
	return 1 // 返回 1 表示成功
}
func luaEditMessage(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msgID := LuaVM.CheckString(2)
	newText := LuaVM.CheckString(3)
	EditMessage(ctx, msgID, newText)
	return 1 // 返回 1 表示成功
}
func luaRecallMessage(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msgID := LuaVM.CheckString(2)
	RecallMessage(ctx, msgID)
	return 1 // 返回 1 表示成功
}
func luaSendToGroupNotice(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	content := LuaVM.CheckString(3)
	SendToGroupNotice(ctx, groupID, content)
	return 1 // 返回 1 表示成功
}
func luaSendLike(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	userID := LuaVM.CheckString(2)
	times := LuaVM.CheckNumber(3)
	SendLike(ctx, userID, int(times))
	return 1 // 返回 1 表示成功
}
func luaSetGroupAdmin(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	userID := LuaVM.CheckString(3)
	enable := LuaVM.CheckBool(4)
	SetGroupAdmin(ctx, groupID, userID, enable)
	return 1 // 返回 1 表示成功
}
func luaSetGroupName(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	groupID := LuaVM.CheckString(2)
	name := LuaVM.CheckString(3)
	SetGroupName(ctx, groupID, name)
	return 1 // 返回 1 表示成功
}
func luaSetSelfLongNick(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	nick := LuaVM.CheckString(2)
	SetSelfLongNick(ctx, nick)
	return 1 // 返回 1 表示成功
}
func luaDiceFormat(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	res := DiceFormat(ctx, s)
	LuaVM.Push(lua.LString(res))
	return 1 // 返回 1 表示成功
}

func luaDiceFormatTmpl(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	s := LuaVM.CheckString(2)
	res := DiceFormatTmpl(ctx, s)
	LuaVM.Push(lua.LString(res))
	return 1 // 返回 1 表示成功
}

func luaShikiSendMsg(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(2).Value.(*Message)
	text := LuaVM.CheckString(3)
	magFromGroup := LuaVM.CheckString(4)
	magFromQQ := LuaVM.CheckString(5)
	if magFromQQ == "" {
		magFromQQ = ctx.Player.UserID
	}
	if magFromGroup != "" && strings.HasPrefix(magFromGroup, "QQ-Group:") == false {
		magFromGroup = fmt.Sprintf("%s%s", "QQ-Group:", magFromGroup)
	}
	if strings.HasPrefix(magFromQQ, "QQ:") == false {
		magFromQQ = fmt.Sprintf("%s%s", "QQ:", magFromQQ)
	}
	msg.Sender.UserID = magFromQQ
	msg.GroupID = magFromGroup
	if magFromGroup == "" {
		ctx.IsPrivate = true
		msg.MessageType = "private"
		msg.Time = int64(time.Now().Unix())
		ctx.Group, ctx.Player = GetPlayerInfoBySender(ctx, msg)
		ReplyPerson(ctx, msg, text)
	} else {
		msg.Time = int64(time.Now().Unix())
		ctx.Group, ctx.Player = GetPlayerInfoBySender(ctx, msg)
		ReplyGroup(ctx, msg, text)
	}
	return 1 // 返回 1 表示成功
}

func luaShikiSetGroupConf(LuaVM *lua.LState) int {
	groupID := LuaVM.CheckString(1)
	key := LuaVM.CheckString(2)
	dice := LuaVM.CheckUserData(3).Value.(*Dice)
	value := LuaVM.Get(4)
	luaType := ""
	if groupID == "" {
		return 0 // 返回 0 表示失败
	} else if strings.HasPrefix(groupID, "QQ-Group:") == false {
		groupID = fmt.Sprintf("%s%s", "QQ-Group:", groupID)
	}
	if value.Type() == lua.LTString {
		luaType = "string"
	} else if value.Type() == lua.LTNumber {
		luaType = "number"
	} else if value.Type() == lua.LTBool {
		luaType = "bool"
	} else {
		return 0 // 返回 0 表示失败
	}
	dice.shikiSetGroupConfig(groupID, key, luaType, value)
	return 1 // 返回 1 表示成功
}

func luaShikiGetGroupConf(LuaVM *lua.LState) int {
	groupID := LuaVM.CheckString(1)
	key := LuaVM.CheckString(2)
	dice := LuaVM.CheckUserData(3).Value.(*Dice)
	if groupID == "" {
		return 0 // 返回 0 表示失败
	} else if strings.HasPrefix(groupID, "QQ-Group:") == false {
		groupID = fmt.Sprintf("%s%s", "QQ-Group:", groupID)
	}
	Type, value := dice.shikiGetGroupConfig(groupID, key)
	if Type == "" {
		return 0 // 返回 0 表示失败
	}
	if Type == "string" {
		LuaVM.Push(lua.LString(value.(string)))
	} else if Type == "number" {
		LuaVM.Push(lua.LNumber(value.(float64)))
	} else if Type == "bool" {
		LuaVM.Push(lua.LBool(value.(bool)))
	}

	return 1 // 返回 1 表示成功
}

func luaShikiSetUserConf(LuaVM *lua.LState) int {
	userID := LuaVM.CheckString(1)
	key := LuaVM.CheckString(2)
	dice := LuaVM.CheckUserData(3).Value.(*Dice)
	value := LuaVM.Get(4)
	luaType := ""
	if userID == "" {
		return 0 // 返回 0 表示失败
	} else if strings.HasPrefix(userID, "QQ:") == false {
		userID = fmt.Sprintf("%s%s", "QQ:", userID)
	}
	if value.Type() == lua.LTString {
		luaType = "string"
	} else if value.Type() == lua.LTNumber {
		luaType = "number"
	} else if value.Type() == lua.LTBool {
		luaType = "bool"
	} else {
		return 0 // 返回 0 表示失败
	}
	dice.shikiSetUserConfig(userID, key, luaType, value)
	return 1 // 返回 1 表示成功
}

func luaShikiGetUserConf(LuaVM *lua.LState) int {
	userID := LuaVM.CheckString(1)
	key := LuaVM.CheckString(2)
	dice := LuaVM.CheckUserData(3).Value.(*Dice)
	if userID == "" {
		return 0 // 返回 0 表示失败
	} else if strings.HasPrefix(userID, "QQ:") == false {
		userID = fmt.Sprintf("%s%s", "QQ:", userID)
	}
	Type, value := dice.shikiGetUserConfig(userID, key)
	if Type == "" {
		return 0 // 返回 0 表示失败
	}
	if Type == "string" {
		LuaVM.Push(lua.LString(value.(string)))
	} else if Type == "number" {
		LuaVM.Push(lua.LNumber(value.(float64)))
	} else if Type == "bool" {
		LuaVM.Push(lua.LBool(value.(bool)))
	}

	return 1 // 返回 1 表示成功
}

func luaShikiGetDiceID(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	LuaVM.Push(lua.LString(ctx.EndPoint.ID))
	return 1 // 返回 1 表示成功
}

func luaShikiLoadLuaFile(LuaVM *lua.LState) int {
	path := LuaVM.CheckString(1)
	d := LuaVM.CheckUserData(2).Value.(*Dice)
	ctx := LuaVM.CheckUserData(3).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(4).Value.(*Message)
	cmdArgs := LuaVM.CheckUserData(5).Value.(*CmdArgs)

	L := lua.NewState()
	defer L.Close()

	//初始化lua全局变量
	LuaVarInit(L, d, ctx, msg, cmdArgs)
	//初始化lua全局函数
	LuaFuncInit(L)

	if err := L.DoFile(path); err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 0
	}

	return 1 // 返回 1 表示成功
}

func luaShikiLoadLuaString(LuaVM *lua.LState) int {
	code := LuaVM.CheckString(1)
	d := LuaVM.CheckUserData(2).Value.(*Dice)
	ctx := LuaVM.CheckUserData(3).Value.(*MsgContext)
	msg := LuaVM.CheckUserData(4).Value.(*Message)
	cmdArgs := LuaVM.CheckUserData(5).Value.(*CmdArgs)

	L := lua.NewState()
	defer L.Close()

	//初始化lua全局变量
	LuaVarInit(L, d, ctx, msg, cmdArgs)
	//初始化lua全局函数
	LuaFuncInit(L)

	if err := L.DoString(code); err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 0
	}

	return 1 // 返回 1 表示成功
}

func luaShikiHTTPGet(LuaVM *lua.LState) int {
	url := LuaVM.CheckString(1)
	resp, err := http.Get(url)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	LuaVM.Push(lua.LString(body))
	return 1 // 返回 1 表示成功
}

func luaShikiHTTPPost(LuaVM *lua.LState) int {
	url := LuaVM.CheckString(1)
	contentType := LuaVM.CheckString(2)
	body := LuaVM.CheckString(3)
	resp, err := http.Post(url, contentType, strings.NewReader(body))
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	LuaVM.Push(lua.LString(respBody))
	return 1 // 返回 1 表示成功
}

func luaShikiHTTPRequest(LuaVM *lua.LState) int {
	method := LuaVM.CheckString(1)
	url := LuaVM.CheckString(2)
	contentType := LuaVM.CheckString(3)
	body := LuaVM.CheckString(4)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	LuaVM.Push(lua.LString(respBody))
	return 1 // 返回 1 表示成功
}

// ----------------------------------------------------------------
func luaDreamJSONEncode(LuaVM *lua.LState) int {
	lv := LuaVM.CheckTable(1)
	goMap := toGoMap(lv)
	jsonData, err := json.Marshal(goMap)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}

	LuaVM.Push(lua.LString(jsonData))
	return 1 // 返回 1 表示成功
}

// luaJSONDecode decodes a JSON string into a Lua table.
func luaDreamJSONDecode(LuaVM *lua.LState) int {
	jsonStr := LuaVM.CheckString(1)

	var goMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &goMap)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}

	luaTable := toLuaTable(LuaVM, goMap)
	LuaVM.Push(luaTable)
	return 1 // 返回 1 表示成功
}

// toGoMap converts a Lua table to a Go map.
func toGoMap(lv *lua.LTable) map[string]interface{} {
	goMap := make(map[string]interface{})
	lv.ForEach(func(key lua.LValue, value lua.LValue) {
		goMap[key.String()] = toGoValue(value)
	})
	return goMap
}

// toGoValue converts a Lua value to a Go value.
func toGoValue(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTString:
		return lv.String()
	case lua.LTNumber:
		return float64(lua.LVAsNumber(lv))
	case lua.LTBool:
		return lua.LVAsBool(lv)
	case lua.LTTable:
		return toGoMap(lv.(*lua.LTable))
	default:
		return nil
	}
}

// toLuaTable converts a Go map to a Lua table.
func toLuaTable(LuaVM *lua.LState, goMap map[string]interface{}) *lua.LTable {
	luaTable := LuaVM.NewTable()
	for key, value := range goMap {
		luaTable.RawSetString(key, toLuaValue(LuaVM, value))
	}
	return luaTable
}

// toLuaValue converts a Go value to a Lua value.
func toLuaValue(LuaVM *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case map[string]interface{}:
		return toLuaTable(LuaVM, v)
	default:
		return lua.LNil
	}
}

// String sub function
func luaDreamStringSub(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	start := LuaVM.CheckInt(2)
	end := LuaVM.CheckInt(3)
	LuaVM.Push(lua.LString(string([]rune(str)[start-1 : end])))
	return 1
}

// String part function
func luaDreamStringPart(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	sep := LuaVM.CheckString(2)
	parts := strings.Split(str, sep)
	table := LuaVM.NewTable()
	for i, part := range parts {
		table.RawSetInt(i+1, lua.LString(part))
	}
	LuaVM.Push(table)
	return 1
}

// String find function
func luaDreamStringFind(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	substr := LuaVM.CheckString(2)
	count := strings.Count(str, substr)
	LuaVM.Push(lua.LNumber(count))
	return 1
}

// String toTable function
func luaDreamStringToTable(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	runes := []rune(str)
	table := LuaVM.NewTable()
	for i, r := range runes {
		table.RawSetInt(i+1, lua.LString(string(r)))
	}
	LuaVM.Push(table)
	return 1
}

// String len function
func luaDreamStringLen(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	length := len([]rune(str))
	LuaVM.Push(lua.LNumber(length))
	return 1
}

// String format function
func luaDreamStringFormat(LuaVM *lua.LState) int {
	str := LuaVM.CheckString(1)
	tab := LuaVM.CheckTable(2)
	vars := make(map[string]string)
	tab.ForEach(func(key, value lua.LValue) {
		vars[key.String()] = value.String()
	})
	for k, v := range vars {
		str = strings.ReplaceAll(str, "{"+k+"}", v)
	}
	LuaVM.Push(lua.LString(str))
	return 1
}

func luaDreamTableType(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	isArray := true
	isObject := true

	table.ForEach(func(key, value lua.LValue) {
		if key.Type() == lua.LTNumber {
			isObject = false
		} else if key.Type() == lua.LTString {
			isArray = false
		}
	})

	if isArray {
		LuaVM.Push(lua.LString("array"))
	} else if isObject {
		LuaVM.Push(lua.LString("object"))
	} else {
		LuaVM.Push(lua.LNil)
	}
	return 1
}

// Function to make a table orderly
func luaDreamTableOrderly(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	newTable := LuaVM.NewTable()

	table.ForEach(func(key, value lua.LValue) {
		if key.Type() == lua.LTNumber {
			newTable.Append(value)
		}
	})
	LuaVM.Push(newTable)
	return 1
}
func luaDreamtableToString(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	result, err := tableToStringInDreamTable(table, "", make(map[*lua.LTable]bool))
	if err != nil {
		LuaVM.RaiseError(err.Error())
		return 0
	}
	LuaVM.Push(lua.LString(result))
	return 1
}

func tableToStringInDreamTable(table *lua.LTable, indent string, visited map[*lua.LTable]bool) (string, error) {
	if visited[table] {
		return "", fmt.Errorf("circular references")
	}
	visited[table] = true

	var sb strings.Builder
	sb.WriteString("{")
	newIndent := indent + "  "

	table.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := luaValueToStringInDreamTable(key)
		valueStr := ""
		if value.Type() == lua.LTTable {
			if visited[value.(*lua.LTable)] {
				valueStr = "circular reference"
			} else {
				var err error
				valueStr, err = tableToStringInDreamTable(value.(*lua.LTable), newIndent, visited)
				if err != nil {
					valueStr = "error"
				}
			}
		} else {
			valueStr = luaValueToStringInDreamTable(value)
		}
		sb.WriteString(fmt.Sprintf("\n%s[%s] -> %s", newIndent, keyStr, valueStr))
	})

	if sb.String() != "{" {
		sb.WriteString(fmt.Sprintf("\n%s}", indent))
	} else {
		sb.WriteString("}")
	}

	return sb.String(), nil
}

func luaValueToStringInDreamTable(value lua.LValue) string {
	switch value.Type() {
	case lua.LTNil:
		return "nil"
	case lua.LTBool:
		return fmt.Sprintf("%t", lua.LVAsBool(value))
	case lua.LTNumber:
		return fmt.Sprintf("%v", lua.LVAsNumber(value))
	case lua.LTString:
		return fmt.Sprintf("%q", lua.LVAsString(value))
	case lua.LTFunction:
		return "function"
	case lua.LTTable:
		return "table"
	case lua.LTUserData:
		return "userdata"
	default:
		return "unknown"
	}
}

// Function to get the length of a table
func luaDreamTableGetNumber(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	length := 0

	table.ForEach(func(key, value lua.LValue) {
		length++
	})
	LuaVM.Push(lua.LNumber(length))
	return 1
}

// Function to sort a table
func luaDreamTableSort(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	key := LuaVM.OptString(2, "")

	var values []lua.LValue
	table.ForEach(func(_, value lua.LValue) {
		values = append(values, value)
	})

	if key == "" {
		sort.Slice(values, func(i, j int) bool {
			return values[i].(lua.LNumber) > values[j].(lua.LNumber)
		})
	} else {
		sort.Slice(values, func(i, j int) bool {
			return values[i].(*lua.LTable).RawGetString(key).(lua.LNumber) > values[j].(*lua.LTable).RawGetString(key).(lua.LNumber)
		})
	}

	newTable := LuaVM.NewTable()
	for _, value := range values {
		newTable.Append(value)
	}
	LuaVM.Push(newTable)
	return 1
}

// Function to clone a table
func luaDreamTableClone(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	newTable := LuaVM.NewTable()

	table.ForEach(func(key, value lua.LValue) {
		newTable.RawSet(key, value)
	})
	LuaVM.Push(newTable)
	return 1
}

// Function to check if two tables are equal
func luaDreamTableEqual(LuaVM *lua.LState) int {
	table1 := LuaVM.CheckTable(1)
	table2 := LuaVM.CheckTable(2)

	if reflect.DeepEqual(table1, table2) {
		LuaVM.Push(lua.LTrue)
	} else {
		LuaVM.Push(lua.LFalse)
	}
	return 1
}

// Function to replace substrings in an array of strings
func luaDreamTableGsub(LuaVM *lua.LState) int {
	table := LuaVM.CheckTable(1)
	old := LuaVM.CheckString(2)
	new := LuaVM.CheckString(3)

	newTable := LuaVM.NewTable()
	table.ForEach(func(_, value lua.LValue) {
		str := strings.Replace(value.String(), old, new, -1)
		newTable.Append(lua.LString(str))
	})
	LuaVM.Push(newTable)
	return 1
}

func luaDreamTableAdd(LuaVM *lua.LState) int {
	// 创建一个新的表用于存储合并结果
	newTable := LuaVM.NewTable()

	// 获取传入参数的数量
	numArgs := LuaVM.GetTop()

	// 遍历所有传入的参数
	for i := 1; i <= numArgs; i++ {
		// 检查参数是否为表
		table := LuaVM.CheckTable(i)

		// 将当前表的所有键值对添加到新表中
		table.ForEach(func(key, value lua.LValue) {
			newTable.RawSet(key, value)
		})
	}

	// 将合并后的新表压入 Lua 栈
	LuaVM.Push(newTable)

	// 返回结果表的数量
	return 1
}

func luaDreamApiGetDiceQQ(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	DiceQQ, _ := strconv.Atoi(strings.ReplaceAll(ctx.EndPoint.UserID, "QQ:", ""))
	LuaVM.Push(lua.LNumber(DiceQQ))
	return 1
}

// Base64 encode function
func luaDreamBase64Encode(LuaVM *lua.LState) int {
	input := LuaVM.CheckString(1)
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	LuaVM.Push(lua.LString(encoded))
	return 1
}

// Base64 decode function
func luaDreamBase64Decode(LuaVM *lua.LState) int {
	input := LuaVM.CheckString(1)
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		LuaVM.Push(lua.LNil)
		LuaVM.Push(lua.LString(err.Error()))
		return 2
	}
	LuaVM.Push(lua.LString(decoded))
	return 1
}

// MD5 hash function
func luaDreamMd5Hash(LuaVM *lua.LState) int {
	input := LuaVM.CheckString(1)
	hash := md5.New()
	io.WriteString(hash, input)
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	LuaVM.Push(lua.LString(hashed))
	return 1
}

// SHA256 hash function
func luaDreamSha256Hash(LuaVM *lua.LState) int {
	input := LuaVM.CheckString(1)
	hash := sha256.New()
	io.WriteString(hash, input)
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	LuaVM.Push(lua.LString(hashed))
	return 1
}

// BKDRHash function
func luaDreamBKDRHash(LuaVM *lua.LState) int {
	input := LuaVM.CheckString(1)
	seed := LuaVM.CheckInt(2)
	hash := BKDRHashInDreamBKDR(input, seed)
	LuaVM.Push(lua.LString(hash))
	return 1
}

// BKDRHash算法
func BKDRHashInDreamBKDR(s string, seed int) string {
	const seed_a = 131  // 31 131 1313 13131 131313 etc.
	const seed_b = 1313 // 131 1313 13131 131313 etc.
	hash := 0
	for _, c := range s {
		hash = (hash*seed_a + int(c)) % seed_b
	}
	return fmt.Sprintf("%d", hash)
}

//----------------------------------------------------------------

func luaZhaoDiceSDKSystemReload(LuaVM *lua.LState) int {
	ctx := LuaVM.CheckUserData(1).Value.(*MsgContext)
	var dm = ctx.Dice.Parent
	dm.RebootRequestChan <- 1
	return 1
}

func luaZhaoDiceSDKTrim(L *lua.LState) int {
	// 获取第一个参数，并确保它是一个字符串
	str := L.CheckString(1)
	// 去除字符串两端的空白字符
	trimmedStr := strings.TrimSpace(str)
	// 将结果压入 Lua 栈
	L.Push(lua.LString(trimmedStr))
	// 返回结果的数量
	return 1
}

func luaZhaoDiceSDKContains(L *lua.LState) int {
	str1 := L.CheckString(1)
	str2 := L.CheckString(2)
	flg := strings.Contains(str1, str2)
	L.Push(lua.LBool(flg))
	// 返回结果的数量
	return 1
}

func LuaFuncInit(LuaVM *lua.LState) {
	LuaVM.SetGlobal("VarSetValueStr", LuaVM.NewFunction(luaVarSetValueStr))
	LuaVM.SetGlobal("VarSetValueInt", LuaVM.NewFunction(luaVarSetValueInt))
	LuaVM.SetGlobal("VarDelValue", LuaVM.NewFunction(luaVarDelValue))
	LuaVM.SetGlobal("VarGetValueInt", LuaVM.NewFunction(luaVarGetValueInt))
	LuaVM.SetGlobal("VarGetValueStr", LuaVM.NewFunction(luaVarGetValueStr))
	LuaVM.SetGlobal("AddBan", LuaVM.NewFunction(luaAddBan))
	LuaVM.SetGlobal("AddTrust", LuaVM.NewFunction(luaAddTrust))
	LuaVM.SetGlobal("RemoveBan", LuaVM.NewFunction(luaRemoveBan))
	LuaVM.SetGlobal("ReplyGroup", LuaVM.NewFunction(luaReplyGroup))
	LuaVM.SetGlobal("ReplyPerson", LuaVM.NewFunction(luaReplyPerson))
	LuaVM.SetGlobal("ReplyToSender", LuaVM.NewFunction(luaReplyToSender))
	LuaVM.SetGlobal("MemberBan", LuaVM.NewFunction(luaMemberBan))
	LuaVM.SetGlobal("MemberWholeBan", LuaVM.NewFunction(luaMemberWholeBan))
	LuaVM.SetGlobal("MemberUnban", LuaVM.NewFunction(luaMemberUnban))
	LuaVM.SetGlobal("MemberKick", LuaVM.NewFunction(luaMemberKick))
	LuaVM.SetGlobal("EditMessage", LuaVM.NewFunction(luaEditMessage))
	LuaVM.SetGlobal("RecallMessage", LuaVM.NewFunction(luaRecallMessage))
	LuaVM.SetGlobal("SendToGroupNotice", LuaVM.NewFunction(luaSendToGroupNotice))
	LuaVM.SetGlobal("SendLike", LuaVM.NewFunction(luaSendLike))
	LuaVM.SetGlobal("SetGroupAdmin", LuaVM.NewFunction(luaSetGroupAdmin))
	LuaVM.SetGlobal("SetGroupName", LuaVM.NewFunction(luaSetGroupName))
	LuaVM.SetGlobal("SetSelfLongNick", LuaVM.NewFunction(luaSetSelfLongNick))
	LuaVM.SetGlobal("DiceFormat", LuaVM.NewFunction(luaDiceFormat))
	LuaVM.SetGlobal("DiceFormatTmpl", LuaVM.NewFunction(luaDiceFormatTmpl))
	LuaVM.SetGlobal("shikisendMsg", LuaVM.NewFunction(luaShikiSendMsg))
	LuaVM.SetGlobal("shikisetGroupConf", LuaVM.NewFunction(luaShikiSetGroupConf))
	LuaVM.SetGlobal("shikigetGroupConf", LuaVM.NewFunction(luaShikiGetGroupConf))
	LuaVM.SetGlobal("shikisetUserConf", LuaVM.NewFunction(luaShikiSetUserConf))
	LuaVM.SetGlobal("shikigetUserConf", LuaVM.NewFunction(luaShikiGetUserConf))
	LuaVM.SetGlobal("shikigetDiceID", LuaVM.NewFunction(luaShikiGetDiceID))
	LuaVM.SetGlobal("shikiloadLuaFile", LuaVM.NewFunction(luaShikiLoadLuaFile))
	LuaVM.SetGlobal("shikiloadLuaString", LuaVM.NewFunction(luaShikiLoadLuaString))
	LuaVM.SetGlobal("shikiHttpGet", LuaVM.NewFunction(luaShikiHTTPGet))
	LuaVM.SetGlobal("shikiHttpPost", LuaVM.NewFunction(luaShikiHTTPPost))
	LuaVM.SetGlobal("shikiHttpRequest", LuaVM.NewFunction(luaShikiHTTPRequest))

	//----------------------------------------------------------------
	DreamLib := LuaVM.NewTable()
	DreamLib.RawSetString("_VERSION", lua.LString("ver4.9.6(206)"))
	DreamLib.RawSetString("version", lua.LString("Dream by 筑梦师V2.0&乐某人 for Tempest Dice"))
	DreamJson := LuaVM.NewTable()
	DreamString := LuaVM.NewTable()
	DreamTable := LuaVM.NewTable()
	DreamBase64 := LuaVM.NewTable()
	DreamMd5 := LuaVM.NewTable()
	DreamSha256 := LuaVM.NewTable()
	DreamBKDR := LuaVM.NewTable()
	LuaVM.SetField(DreamLib, "json", DreamJson)
	LuaVM.SetField(DreamJson, "encode", LuaVM.NewFunction(luaDreamJSONEncode))
	LuaVM.SetField(DreamJson, "decode", LuaVM.NewFunction(luaDreamJSONDecode))
	LuaVM.SetField(DreamLib, "string", DreamString)
	LuaVM.SetField(DreamString, "sub", LuaVM.NewFunction(luaDreamStringSub))
	LuaVM.SetField(DreamString, "part", LuaVM.NewFunction(luaDreamStringPart))
	LuaVM.SetField(DreamString, "find", LuaVM.NewFunction(luaDreamStringFind))
	LuaVM.SetField(DreamString, "totable", LuaVM.NewFunction(luaDreamStringToTable))
	LuaVM.SetField(DreamString, "len", LuaVM.NewFunction(luaDreamStringLen))
	LuaVM.SetField(DreamString, "format", LuaVM.NewFunction(luaDreamStringFormat))
	LuaVM.SetField(DreamLib, "table", DreamTable)
	LuaVM.SetField(DreamTable, "type", LuaVM.NewFunction(luaDreamTableType))
	LuaVM.SetField(DreamTable, "orderly", LuaVM.NewFunction(luaDreamTableOrderly))
	LuaVM.SetField(DreamTable, "getnumber", LuaVM.NewFunction(luaDreamTableGetNumber))
	LuaVM.SetField(DreamTable, "sort", LuaVM.NewFunction(luaDreamTableSort))
	LuaVM.SetField(DreamTable, "clone", LuaVM.NewFunction(luaDreamTableClone))
	LuaVM.SetField(DreamTable, "equal", LuaVM.NewFunction(luaDreamTableEqual))
	LuaVM.SetField(DreamTable, "gsub", LuaVM.NewFunction(luaDreamTableGsub))
	LuaVM.SetField(DreamTable, "add", LuaVM.NewFunction(luaDreamTableAdd))
	LuaVM.SetField(DreamTable, "tostring", LuaVM.NewFunction(luaDreamtableToString))
	LuaVM.SetField(DreamTable, "tostr", LuaVM.NewFunction(luaDreamtableToString))
	LuaVM.SetField(DreamLib, "base64", DreamBase64)
	LuaVM.SetField(DreamBase64, "encode", LuaVM.NewFunction(luaDreamBase64Encode))
	LuaVM.SetField(DreamBase64, "decode", LuaVM.NewFunction(luaDreamBase64Decode))
	LuaVM.SetField(DreamLib, "md5", DreamMd5)
	LuaVM.SetField(DreamMd5, "hash", LuaVM.NewFunction(luaDreamMd5Hash))
	LuaVM.SetField(DreamLib, "sha256", DreamSha256)
	LuaVM.SetField(DreamSha256, "hash", LuaVM.NewFunction(luaDreamSha256Hash))
	LuaVM.SetField(DreamLib, "BKDR", DreamBKDR)
	LuaVM.SetField(DreamBKDR, "hash", LuaVM.NewFunction(luaDreamBKDRHash))
	LuaVM.SetGlobal("dream", DreamLib)
	//----------------------------------------------------------------
	ZhaoDiceSDK := LuaVM.NewTable()
	ZhaoDiceSDKSystem := LuaVM.NewTable()
	LuaVM.SetField(ZhaoDiceSDK, "trim", LuaVM.NewFunction(luaZhaoDiceSDKTrim))
	LuaVM.SetField(ZhaoDiceSDK, "contains", LuaVM.NewFunction(luaZhaoDiceSDKContains))
	LuaVM.SetField(ZhaoDiceSDK, "system", ZhaoDiceSDKSystem)
	LuaVM.SetField(ZhaoDiceSDKSystem, "reload", LuaVM.NewFunction(luaZhaoDiceSDKSystemReload))
	LuaVM.SetGlobal("zhaodicesdk", ZhaoDiceSDK)
}
