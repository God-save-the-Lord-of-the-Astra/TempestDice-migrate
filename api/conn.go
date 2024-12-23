package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"tempestdice/dice"
)

func ImConnections(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	return c.JSON(http.StatusOK, myDice.ImSession.EndPoints)
}

func ImConnectionsGet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				return c.JSON(http.StatusOK, i)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsSetEnable(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID     string `form:"id" json:"id"`
		Enable bool   `form:"enable" json:"enable"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				i.SetEnable(myDice, v.Enable)
				return c.JSON(http.StatusOK, i)
			}
		}
	}

	myDice.LastUpdatedTime = time.Now().Unix()
	myDice.Save(false)
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsSetData(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID                  string `form:"id" json:"id"`
		Protocol            int    `form:"protocol" json:"protocol"`
		AppVersion          string `form:"appVersion" json:"appVersion"`
		IgnoreFriendRequest bool   `json:"ignoreFriendRequest"` // 忽略好友请求
		UseSignServer       bool   `json:"useSignServer"`
		ExtraArgs           string `json:"extraArgs"`
		SignServerConfig    *dice.SignServerConfig
	}{}

	err := c.Bind(&v)
	if err != nil {
		myDice.Save(false)
		return c.JSON(http.StatusNotFound, nil)
	}
	for _, i := range myDice.ImSession.EndPoints {
		if i.ID != v.ID {
			continue
		}
		ad := i.Adapter.(*dice.PlatformAdapterGocq)
		if i.ProtocolType != "onebot" {
			i.ProtocolType = "onebot"
		}
		ad.SetQQProtocol(v.Protocol)
		ad.InPackGoCqhttpAppVersion = v.AppVersion
		if v.UseSignServer {
			ad.SetSignServer(v.SignServerConfig)
			ad.UseSignServer = v.UseSignServer
			ad.SignServerConfig = v.SignServerConfig
		}
		ad.IgnoreFriendRequest = v.IgnoreFriendRequest
		ad.ExtraArgs = v.ExtraArgs
		return c.JSON(http.StatusOK, i)
	}
	myDice.Save(false)
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsRWSignServerUrl(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID                string `form:"id" json:"id"`
		SignServerUrl     string `form:"signServerUrl" json:"signServerUrl"`
		W                 bool   `form:"w" json:"w"`
		SignServerVersion string `form:"signServerVersion" json:"signServerVersion"`
	}{}

	err := c.Bind(&v)
	if err != nil {
		myDice.Save(false)
		return c.JSON(http.StatusNotFound, nil)
	}
	for _, i := range myDice.ImSession.EndPoints {
		if i.ID != v.ID {
			continue
		}
		if i.ProtocolType == "onebot" {
			pa := i.Adapter.(*dice.PlatformAdapterGocq)
			if pa.BuiltinMode == "lagrange" {
				signServerUrl, signServerVersion := dice.RWLagrangeSignServerUrl(myDice, i, v.SignServerUrl, v.W, v.SignServerVersion)
				if signServerUrl != "" {
					return Success(&c, Response{
						"signServerUrl":     signServerUrl,
						"signServerVersion": signServerVersion,
					})
				}
			}
		}
	}
	return Error(&c, "读取signServerUrl字段失败", Response{})
}

func ImConnectionsDel(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for index, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				// 禁用该endpoint防止出问题
				i.SetEnable(myDice, false)
				// 待删除的EPInfo落库，保留其统计数据
				i.StatsDump(myDice)
				// TODO: 注意 这个好像很不科学
				// i.diceServing = false
				switch i.Platform {
				case "QQ":
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					if i.ProtocolType == "onebot" {
						pa := i.Adapter.(*dice.PlatformAdapterGocq)
						if pa.BuiltinMode == "lagrange" {
							dice.BuiltinQQServeProcessKillBase(myDice, i, true)
							// 经测试，若不延时，可能导致清理对应目录失败（原因：文件被占用）
							time.Sleep(1 * time.Second)
							dice.LagrangeServeRemoveConfig(myDice, i)
						} else {
							dice.BuiltinQQServeProcessKill(myDice, i)
						}
					}
					return c.JSON(http.StatusOK, i)
				case "DISCORD":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "KOOK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "TG":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "MC":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "DODO":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "DINGTALK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "SLACK":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				case "SEALCHAT":
					i.Adapter.SetEnable(false)
					myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints[:index], myDice.ImSession.EndPoints[index+1:]...)
					return c.JSON(http.StatusOK, i)
				}
			}
		}
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsQrcodeGet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	for _, i := range myDice.ImSession.EndPoints {
		if i.ID != v.ID {
			continue
		}
		switch i.ProtocolType {
		case "onebot", "":
			pa := i.Adapter.(*dice.PlatformAdapterGocq)
			if pa.GoCqhttpState == dice.StateCodeInLoginQrCode {
				return c.JSON(http.StatusOK, map[string]string{
					"img": "data:image/png;base64," + base64.StdEncoding.EncodeToString(pa.GoCqhttpQrcodeData),
				})
			}
		}
		return c.JSON(http.StatusOK, i)
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsCaptchaSet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID   string `form:"id" json:"id"`
		Code string `form:"code" json:"code"`
	}{}
	err := c.Bind(&v)
	if err != nil {
		return err
	}

	for _, i := range myDice.ImSession.EndPoints {
		if i.ID == v.ID {
			switch i.ProtocolType {
			case "onebot", "":
				pa := i.Adapter.(*dice.PlatformAdapterGocq)
				if pa.GoCqhttpState == dice.GoCqhttpStateCodeInLoginBar {
					pa.GoCqhttpLoginCaptcha = v.Code
					return c.String(http.StatusOK, "")
				}
			}
		}
	}
	return c.String(http.StatusNotFound, "")
}

func ImConnectionsSmsCodeSet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID   string `form:"id" json:"id"`
		Code string `form:"code" json:"code"`
	}{}
	err := c.Bind(&v)

	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				switch i.ProtocolType {
				case "onebot", "":
					pa := i.Adapter.(*dice.PlatformAdapterGocq)
					if pa.GoCqhttpState == dice.GoCqhttpStateCodeInLoginVerifyCode {
						pa.GoCqhttpLoginVerifyCode = v.Code
						return c.JSON(http.StatusOK, map[string]string{})
					}
				}
				return c.JSON(http.StatusOK, i)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsSmsCodeGet(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)

	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				switch i.ProtocolType {
				case "onebot", "":
					pa := i.Adapter.(*dice.PlatformAdapterGocq)
					return c.JSON(http.StatusOK, map[string]string{"tip": pa.GoCqhttpSmsNumberTip})
				}
				return c.JSON(http.StatusOK, i)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsGocqhttpRelogin(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				fmt.Print("!!! relogin ", v.ID)
				i.Adapter.DoRelogin()
				return c.JSON(http.StatusOK, nil)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsWalleQRelogin(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}

	v := struct {
		ID string `form:"id" json:"id"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		for _, i := range myDice.ImSession.EndPoints {
			if i.ID == v.ID {
				i.Adapter.DoRelogin()
				return c.JSON(http.StatusOK, nil)
			}
		}
	}
	return c.JSON(http.StatusNotFound, nil)
}

func ImConnectionsGocqConfigDownload(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	id := c.QueryParam("id")
	for _, i := range myDice.ImSession.EndPoints {
		if i.ID == id {
			buf := packGocqConfig(i.RelWorkDir)
			return c.Blob(http.StatusOK, "", buf.Bytes())
		}
	}

	return c.String(http.StatusNotFound, "")
}

type AddDiscordEcho struct {
	Token              string
	ProxyURL           string
	ReverseProxyUrl    string
	ReverseProxyCDNUrl string
}

func ImConnectionsAddBuiltinGocq(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account          string                 `yaml:"account" json:"account"`
		Password         string                 `yaml:"password" json:"password"`
		Protocol         int                    `json:"protocol"`
		AppVersion       string                 `json:"appVersion"`
		UseSignServer    bool                   `json:"useSignServer"`
		SignServerConfig *dice.SignServerConfig `json:"signServerConfig"`
		// ConnectUrl        string `yaml:"connectUrl" json:"connectUrl"`               // 连接地址
		// Platform          string `yaml:"platform" json:"platform"`                   // 平台，如QQ、QQ频道
		// Enable            bool   `yaml:"enable" json:"enable"`                       // 是否启用
		// Type              string `yaml:"type" json:"type"`                           // 协议类型，如onebot、koishi等
		// UseInPackGoCqhttp bool   `yaml:"useInPackGoCqhttp" json:"useInPackGoCqhttp"` // 是否使用内置的gocqhttp
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem(v.Account)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.InPackGoCqhttpProtocol = v.Protocol
		pa.InPackGoCqhttpPassword = v.Password
		pa.InPackGoCqhttpAppVersion = v.AppVersion
		pa.Session = myDice.ImSession
		pa.UseSignServer = v.UseSignServer
		pa.SignServerConfig = v.SignServerConfig

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()

		dice.GoCqhttpServe(myDice, conn, dice.GoCqhttpLoginInfo{
			Password:         v.Password,
			Protocol:         v.Protocol,
			AppVersion:       v.AppVersion,
			IsAsyncRun:       true,
			UseSignServer:    v.UseSignServer,
			SignServerConfig: v.SignServerConfig,
		})
		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddGocqSeparate(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account     string `yaml:"account" json:"account"`
		ConnectURL  string `yaml:"connectUrl" json:"connectUrl"`   // 连接地址
		AccessToken string `yaml:"accessToken" json:"accessToken"` // 访问令牌
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem("")
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession

		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.Session = myDice.ImSession

		// 三项设置
		conn.RelWorkDir = "x" // 此选项已无意义
		pa.ConnectURL = v.ConnectURL
		pa.AccessToken = v.AccessToken

		pa.UseInPackClient = false

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		conn.SetEnable(myDice, true)

		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddReverseWs(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"testMode": true,
		})
	}

	v := struct {
		Account     string `yaml:"account" json:"account"`
		ReverseAddr string `yaml:"reverseAddr" json:"reverseAddr"`
	}{}

	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewGoCqhttpConnectInfoItem(v.Account)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession

		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		pa.Session = myDice.ImSession

		pa.IsReverse = true
		pa.ReverseAddr = v.ReverseAddr

		pa.UseInPackClient = false

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		conn.SetEnable(myDice, true)

		myDice.LastUpdatedTime = time.Now().Unix()
		myDice.Save(false)
		return c.JSON(http.StatusOK, conn)
	}
	return c.String(430, "")
}

func ImConnectionsAddBuiltinLagrange(c echo.Context) error {
	if !doAuth(c) {
		return c.JSON(http.StatusForbidden, nil)
	}
	if dm.JustForTest {
		return c.JSON(http.StatusOK, Response{"testMode": true})
	}

	v := struct {
		Account           string `yaml:"account" json:"account"`
		SignServerUrl     string `yaml:"signServerUrl" json:"signServerUrl"`
		SignServerVersion string `yaml:"signServerVersion" json:"signServerVersion"`
	}{}
	err := c.Bind(&v)
	if err == nil {
		uid := v.Account
		if checkUidExists(c, uid) {
			return nil
		}

		conn := dice.NewLagrangeConnectInfoItem(v.Account)
		conn.UserID = dice.FormatDiceIDQQ(uid)
		conn.Session = myDice.ImSession
		pa := conn.Adapter.(*dice.PlatformAdapterGocq)
		// pa.InPackGoCqhttpProtocol = v.Protocol
		pa.Session = myDice.ImSession

		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
		myDice.LastUpdatedTime = time.Now().Unix()
		uin, err := strconv.ParseInt(v.Account, 10, 64)
		if err != nil {
			return err
		}
		dice.LagrangeServe(myDice, conn, dice.LagrangeLoginInfo{
			UIN:               uin,
			SignServerUrl:     v.SignServerUrl,
			SignServerVersion: v.SignServerVersion,
			IsAsyncRun:        true,
		})
		return c.JSON(http.StatusOK, v)
	}

	return c.String(430, "")
}

// func ImConnectionsAddLagrangeGO(c echo.Context) error {
//	if !doAuth(c) {
//		return c.JSON(http.StatusForbidden, nil)
//	}
//	if dm.JustForTest {
//		return Success(&c, Response{"testMode": true})
//	}
//
//	v := struct {
//		Account       string `yaml:"account" json:"account"`
//		CustomSignUrl string `yaml:"signServerUrl" json:"signServerUrl"`
//	}{}
//	err := c.Bind(&v)
//	if err == nil {
//		uid := v.Account
//		if checkUidExists(c, uid) {
//			return nil
//		}
//		uin, err := strconv.ParseInt(v.Account, 10, 64)
//		if err != nil {
//			return err
//		}
//		conn := dice.NewLagrangeGoConnItem(uint32(uin), v.CustomSignUrl)
//		conn.UserID = dice.FormatDiceIDQQ(uid)
//		conn.Session = myDice.ImSession
//		pa := conn.Adapter.(*dice.PlatformAdapterLagrangeGo)
//		pa.Session = myDice.ImSession
//
//		myDice.ImSession.EndPoints = append(myDice.ImSession.EndPoints, conn)
//		myDice.LastUpdatedTime = time.Now().Unix()
//
//		dice.ServeLagrangeGo(myDice, conn)
//		return c.JSON(http.StatusOK, v)
//	}
//
//	return c.String(430, "")
// }
