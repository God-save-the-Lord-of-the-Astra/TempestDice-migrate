package dice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RegisterBuiltinCloudCommands(d *Dice) {
	helpForShikiCloudBlack := ".cloud sync //与云黑服务器手动同步一次\n" +
		".cloud autosync //与云黑服务器每天自动同步一次(这是个饼)\n" +
		".cloud server list //查看云黑服务器列表\n" +
		".cloud server +/- <名称> <链接> [<权重>] //添加/删除云黑服务器\n" +
		".cloud server add/del <名称> <链接> [<权重>] //添加/删除云黑服务器"

	cmdCloudBlack := &CmdItemInfo{
		Name:      "cloud",
		ShortHelp: helpForShikiCloudBlack,
		Help:      "同步云黑指令:\n" + helpForShikiCloudBlack,
		Solve: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) CmdExecuteResult {
			cmdArgs.ChopPrefixToArgsWith("sync", "autosync")
			if ctx.PrivilegeLevel < 100 {
				ReplyToSender(ctx, msg, DiceFormatTmpl(ctx, "核心:提示_无权限_非master"))
				return CmdExecuteResult{Matched: true, Solved: true}
			}

			ctx.EndPoint.Platform = "QQ"

			type blackunit struct {
				BlackQQ      string
				BlackGroup   string
				WarningID    string
				BlackComment string
				ErasedStatus bool
				ServerWeight int
			}

			type jsonElement struct {
				Wid       int    `json:"wid"`
				FromGroup int    `json:"fromGroup"`
				FromQQ    int    `json:"fromQQ"`
				Type      string `json:"type"`
				Note      string `json:"note"`
				IsErased  int    `json:"isErased"`
			}

			fetchAndParseJSON_shikiCloudBlack := func(url string, serverWeight int) ([]blackunit, error) {
				// 发送 HTTP GET 请求
				resp, err := http.Get(url)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()

				// 读取响应体
				body, _ := io.ReadAll(resp.Body)
				var jsonData []jsonElement
				err = json.Unmarshal(body, &jsonData)
				if err != nil {
					return nil, err
				}

				// 将 JSON 数据转换为 blackunit 结构体数组
				var blackUnits []blackunit
				for _, item := range jsonData {
					unit := blackunit{
						BlackQQ:      strconv.Itoa(item.FromQQ),
						BlackGroup:   strconv.Itoa(item.FromGroup),
						WarningID:    strconv.Itoa(item.Wid),
						BlackComment: item.Type + " " + item.Note,
						ErasedStatus: item.IsErased != 0,
						ServerWeight: serverWeight,
					}
					blackUnits = append(blackUnits, unit)
				}

				return blackUnits, nil
			}

			val := strings.ToLower(cmdArgs.GetArgN(1))
			subval := strings.ToLower(cmdArgs.GetArgN(2))
			switch val {
			case "sync":
				ReplyToSender(ctx, msg, "正在同步云黑...")
				time.Sleep(1000 * time.Millisecond)

				type BlacklistData struct {
					BlackUnits   []blackunit
					ServerWeight int
				}

				var allBlackUnits []BlacklistData
				ServerListLen := len(d.BlackServerList)
				AliveListLen := 0
				if ServerListLen == 0 {
					ReplyToSender(ctx, msg, "云黑服务器列表为空")
					return CmdExecuteResult{Matched: true, Solved: true}
				}

				// 获取所有服务器的黑名单和擦除名单
				for _, server := range d.BlackServerList {
					url := server.ServerUrl
					weight := server.ServerWeight
					blackUnits, err := fetchAndParseJSON_shikiCloudBlack(url, weight)
					if err != nil {
						ReplyToSender(ctx, msg, fmt.Sprintf("从服务器 %s 获取数据失败: %v", server.ServerName, err))
						continue
					}

					blackUnitsData := BlacklistData{
						BlackUnits:   blackUnits,
						ServerWeight: server.ServerWeight,
					}
					allBlackUnits = append(allBlackUnits, blackUnitsData)
					AliveListLen++
				}
				time.Sleep(1000 * time.Millisecond)
				ReplyToSender(ctx, msg, fmt.Sprintf("%s%d%s%d%s", "列表中共有: ", ServerListLen, "组服务器，本次同步其中: ", AliveListLen, "组，接下来开始同步，请耐心等待。"))
				// 合并黑名单和擦除名单
				mergedBlackUnits := make(map[string]blackunit)
				mergedErasedUnits := make(map[string]blackunit)

				// 遍历所有服务器返回的黑名单数据
				for _, data := range allBlackUnits {
					// 遍历每个服务器返回的黑名单条目
					for _, blackitem := range data.BlackUnits {
						if blackitem.ErasedStatus {
							// 如果条目是擦除状态
							if existingBlackItem, exists := mergedBlackUnits[blackitem.BlackQQ]; exists {
								// 如果已经存在相同的黑名单条目，比较权重
								if existingBlackItem.ServerWeight < blackitem.ServerWeight {
									// 如果擦除条目的权重大于黑名单条目，更新为擦除条目
									delete(mergedBlackUnits, blackitem.BlackQQ)
									mergedErasedUnits[blackitem.BlackQQ] = blackitem
								}
							} else {
								// 如果不存在相同的黑名单条目，直接添加擦除条目
								if existingErasedItem, exists := mergedErasedUnits[blackitem.BlackQQ]; exists {
									// 如果已经存在相同的擦除条目，比较权重
									if existingErasedItem.ServerWeight < blackitem.ServerWeight {
										mergedErasedUnits[blackitem.BlackQQ] = blackitem
									}
								} else {
									mergedErasedUnits[blackitem.BlackQQ] = blackitem
								}
							}
						} else {
							// 如果条目是黑名单状态
							if existingErasedItem, exists := mergedErasedUnits[blackitem.BlackQQ]; exists {
								// 如果已经存在相同的擦除条目，比较权重
								if existingErasedItem.ServerWeight <= blackitem.ServerWeight {
									// 如果黑名单条目的权重大于或等于擦除条目，更新为黑名单条目
									delete(mergedErasedUnits, blackitem.BlackQQ)
									mergedBlackUnits[blackitem.BlackQQ] = blackitem
								}
							} else {
								// 如果不存在相同的擦除条目，直接添加黑名单条目
								if existingBlackItem, exists := mergedBlackUnits[blackitem.BlackQQ]; exists {
									// 如果已经存在相同的黑名单条目，比较权重
									if existingBlackItem.ServerWeight < blackitem.ServerWeight {
										mergedBlackUnits[blackitem.BlackQQ] = blackitem
									}
								} else {
									mergedBlackUnits[blackitem.BlackQQ] = blackitem
								}
							}
						}
					}
				}

				blackGroupCnt := 0
				blackGroupNewCnt := 0
				blackQQCnt := 0
				blackQQNewCnt := 0
				erasedCnt := 0
				erasedNewCnt := 0

				// 根据合并后的名单进行同步操作
				for _, blackitem := range mergedBlackUnits {
					qqTobeBlack := FormatDiceID(ctx, blackitem.BlackQQ, false)
					groupTobeBlack := FormatDiceID(ctx, blackitem.BlackGroup, true)
					if !blackitem.ErasedStatus {
						item, ok := d.BanList.GetByID(qqTobeBlack)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(qqTobeBlack, d.BanList.ThresholdBan, "云黑api", blackitem.BlackComment, ctx)
							blackQQNewCnt++
						}
						blackQQCnt++
						item, ok = d.BanList.GetByID(groupTobeBlack)
						if !ok || (item.Rank != BanRankBanned && item.Rank != BanRankTrusted && item.Rank != BanRankWarn) {
							d.BanList.AddScoreBase(groupTobeBlack, d.BanList.ThresholdBan, "云黑api", blackitem.BlackComment, ctx)
							blackGroupNewCnt++
						}
						blackGroupCnt++
					}
				}

				for _, erasedItem := range mergedErasedUnits {
					qqTobeBlack := FormatDiceID(ctx, erasedItem.BlackQQ, false)
					groupTobeBlack := FormatDiceID(ctx, erasedItem.BlackGroup, true)
					erasedCnt++
					item, ok := d.BanList.GetByID(qqTobeBlack)
					if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
						item.Score = 0
						item.Rank = BanRankNormal
						erasedNewCnt++
					}
					item, ok = d.BanList.GetByID(groupTobeBlack)
					if ok && (item.Rank == BanRankBanned || item.Rank == BanRankTrusted || item.Rank == BanRankWarn) {
						item.Score = 0
						item.Rank = BanRankNormal
					}
				}
				time.Sleep(1000 * time.Millisecond)
				ReplyToSender(ctx, msg, fmt.Sprintf(
					"共计从云黑api获取黑名单群组:%d个，新增:%d个；黑名单用户:%d名，新增:%d名。并有%d组已在云端消除黑名单记录，新增%d组✓",
					blackGroupCnt, blackGroupNewCnt, blackQQCnt, blackQQNewCnt, erasedCnt, erasedNewCnt,
				))

				return CmdExecuteResult{Matched: true, Solved: true}

			case "server":
				switch subval {
				case "list":
					text := ""
					for _, s := range d.BlackServerList {
						text = fmt.Sprintf("%s\t-%s  %d\n", text, s.ServerName, s.ServerWeight)
					}
					if text == "" {
						text = "\t-无\t\t\t无"
					}
					reply := "云黑服务器列表: \n\t-名称\t权重\n" + text
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				case "add":
					newServerName := cmdArgs.GetArgN(3)
					newServerUrl := cmdArgs.GetArgN(4)
					newsw := cmdArgs.GetArgN(5)
					if newServerName == "" {
						ReplyToSender(ctx, msg, "请写出云黑服务器名称")
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					if newServerUrl == "" {
						ReplyToSender(ctx, msg, "请写出云黑服务器地址")
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					if strings.ToLower(newServerName) == "default" {
						newDefaultServerItem := BlackServerListWithWeight{
							ServerName:   "溯洄云黑",
							ServerUrl:    "https://shiki.stringempty.xyz/blacklist/checked.json?",
							ServerWeight: 1,
						}
						d.BlackServerList = append(d.BlackServerList, newDefaultServerItem)
						reply := fmt.Sprintf("%s%s%s%s%s%d%s", "成功添加默认云黑服务器: ", newDefaultServerItem.ServerName, "\n服务器地址: ", newDefaultServerItem.ServerUrl, "\n服务器权重: ", newDefaultServerItem.ServerWeight, "✓")
						d.Save(false)
						ReplyToSender(ctx, msg, reply)
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					newServerWeight := 1
					if newsw == "" {
						newServerWeight = 1
					} else {
						newServerWeight, _ = strconv.Atoi(newsw)
					}
					for index, item := range d.BlackServerList {
						if item.ServerName == newServerName {
							d.BlackServerList[index].ServerUrl = newServerUrl
							d.BlackServerList[index].ServerWeight = newServerWeight
							reply := fmt.Sprintf("%s%s%s%s%s%d%s", "成功编辑云黑服务器: ", newServerName, "\n服务器地址: ", newServerUrl, "\n服务器权重: ", newServerWeight, "✓")
							ReplyToSender(ctx, msg, reply)
							return CmdExecuteResult{Matched: true, Solved: true}
						}
					}
					newServerItem := BlackServerListWithWeight{
						ServerName:   newServerName,
						ServerUrl:    newServerUrl,
						ServerWeight: newServerWeight,
					}
					d.BlackServerList = append(d.BlackServerList, newServerItem)
					reply := fmt.Sprintf("%s%s%s%s%s%d%s", "成功添加云黑服务器: ", newServerName, "\n服务器地址: ", newServerUrl, "\n服务器权重: ", newServerWeight, "✓")
					d.Save(false)
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				case "+":
					newServerName := cmdArgs.GetArgN(3)
					newServerUrl := cmdArgs.GetArgN(4)
					newsw := cmdArgs.GetArgN(5)
					if newServerName == "" {
						ReplyToSender(ctx, msg, "请写出云黑服务器名称")
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					if newServerUrl == "" {
						ReplyToSender(ctx, msg, "请写出云黑服务器地址")
						return CmdExecuteResult{Matched: true, Solved: true}
					}
					newServerWeight := 1
					if newsw == "" {
						newServerWeight = 1
					} else {
						newServerWeight, _ = strconv.Atoi(newsw)
					}
					for index, item := range d.BlackServerList {
						if item.ServerName == newServerName {
							d.BlackServerList[index].ServerUrl = newServerUrl
							d.BlackServerList[index].ServerWeight = newServerWeight
							reply := fmt.Sprintf("%s%s%s%s%s%d%s", "成功编辑云黑服务器: ", newServerName, "\n服务器地址: ", newServerUrl, "\n服务器权重: ", newServerWeight, "✓")
							ReplyToSender(ctx, msg, reply)
							return CmdExecuteResult{Matched: true, Solved: true}
						}
					}
					newServerItem := BlackServerListWithWeight{
						ServerName:   newServerName,
						ServerUrl:    newServerUrl,
						ServerWeight: newServerWeight,
					}
					d.BlackServerList = append(d.BlackServerList, newServerItem)
					reply := fmt.Sprintf("%s%s%s%s%s%d%s", "成功添加云黑服务器: ", newServerName, "\n服务器地址: ", newServerUrl, "\n服务器权重: ", newServerWeight, "✓")
					d.Save(false)
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				case "del":
					newServerElement := cmdArgs.GetArgN(3)
					reply := ""
					for index, item := range d.BlackServerList {
						if item.ServerName == newServerElement || item.ServerUrl == newServerElement {
							delServerName := item.ServerName
							delServerUrl := item.ServerUrl
							delServerWeight := item.ServerWeight
							d.BlackServerList = append(d.BlackServerList[:index], d.BlackServerList[index+1:]...)
							reply = fmt.Sprintf("%s%s%s%s%s%d%s", "成功删除云黑服务器: ", delServerName, "\n服务器地址: ", delServerUrl, "\n服务器权重: ", delServerWeight, "✓")
						}
					}
					if reply == "" {
						reply = "没有找到指定云黑服务器，请先添加。"
					}
					d.Save(false)
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}

				case "-":
					newServerElement := cmdArgs.GetArgN(3)
					reply := ""
					for index, item := range d.BlackServerList {
						if item.ServerName == newServerElement || item.ServerUrl == newServerElement {
							delServerName := item.ServerName
							delServerUrl := item.ServerUrl
							delServerWeight := item.ServerWeight
							d.BlackServerList = append(d.BlackServerList[:index], d.BlackServerList[index+1:]...)
							reply = fmt.Sprintf("%s%s%s%s%s%d%s", "成功删除云黑服务器: ", delServerName, "\n服务器地址: ", delServerUrl, "\n服务器权重: ", delServerWeight, "✓")
						}
					}
					if reply == "" {
						reply = "没有找到指定云黑服务器，请先添加。"
					}
					d.Save(false)
					ReplyToSender(ctx, msg, reply)
					return CmdExecuteResult{Matched: true, Solved: true}
				default:
					return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
				}
			default:
				return CmdExecuteResult{Matched: true, Solved: true, ShowHelp: true}
			}
		},
	}
	d.RegisterExtension(&ExtInfo{
		Name:            "cloud", // 扩展的名称，需要用于指令中，写简短点      2024.05.10: 目前被看成是 function 的缩写了（
		Version:         "1.0.0",
		Brief:           "云黑同步指令",
		AutoActive:      true, // 是否自动开启
		ActiveOnPrivate: true,
		Author:          "海棠,星界之主",
		Official:        true,
		OnCommandReceived: func(ctx *MsgContext, msg *Message, cmdArgs *CmdArgs) {
		},
		OnLoad: func() {
		},
		GetDescText: GetExtensionDesc,
		CmdMap: CmdMapCls{
			"cloud": cmdCloudBlack,
		},
	})
}
