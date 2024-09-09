package main

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	//"github.com/polevpn/webview"
	// _ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap/zapcore"

	"sealdice-core/api"
	"sealdice-core/dice"
	diceLogger "sealdice-core/dice/logger"
	"sealdice-core/dice/model"
	"sealdice-core/migrate"
	"sealdice-core/static"
)

/**
二进制目录结构:
data/configs
data/extensions
data/logs

extensions/
*/

func cleanupCreate(diceManager *dice.DiceManager) func() {
	return func() {
		logger.Info("程序即将退出，进行清理……")
		err := recover()
		if err != nil {
			showWindow()
			logger.Errorf("异常: %v\n堆栈: %v", err, string(debug.Stack()))
			// 顺便修正一下上面这个，应该是木落忘了。
			if runtime.GOOS == "windows" {
				exec.Command("pause") // windows专属
			}
		}

		if !diceManager.CleanupFlag.CompareAndSwap(0, 1) {
			// 尝试更新cleanup标记，如果已经为1则退出
			return
		}

		for _, i := range diceManager.Dice {
			if i.IsAlreadyLoadConfig {
				i.BanList.SaveChanged(i)
				i.Save(true)
				for _, j := range i.ExtList {
					if j.Storage != nil {
						// 关闭
						err := j.StorageClose()
						if err != nil {
							showWindow()
							logger.Errorf("异常: %v\n堆栈: %v", err, string(debug.Stack()))
							// 木落没有加该检查 补充上
							if runtime.GOOS == "windows" {
								exec.Command("pause") // windows专属
							}
						}
					}
				}
				i.IsAlreadyLoadConfig = false
			}
		}

		for _, i := range diceManager.Dice {
			d := i
			(func() {
				defer func() {
					_ = recover()
				}()
				dbData := d.DBData
				if dbData != nil {
					d.DBData = nil
					_ = dbData.Close()
				}
			})()

			(func() {
				defer func() {
					_ = recover()
				}()
				dbLogs := d.DBLogs
				if dbLogs != nil {
					d.DBLogs = nil
					_ = dbLogs.Close()
				}
			})()

			(func() {
				defer func() {
					_ = recover()
				}()
				cm := d.CensorManager
				if cm != nil && cm.DB != nil {
					dbCensor := cm.DB
					cm.DB = nil
					_ = dbCensor.Close()
				}
			})()
		}

		// 清理gocqhttp
		for _, i := range diceManager.Dice {
			if i.ImSession != nil && i.ImSession.EndPoints != nil {
				for _, j := range i.ImSession.EndPoints {
					dice.BuiltinQQServeProcessKill(i, j)
				}
			}
		}

		if diceManager.Help != nil {
			diceManager.Help.Close()
		}
		if diceManager.IsReady {
			diceManager.Save()
		}
		if diceManager.Cron != nil {
			diceManager.Cron.Stop()
		}
	}
}

/*
func fixTimezone() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}
*/

func main() {
	initStartTime := time.Now().UnixMicro()
	var opts struct {
		Version                bool   `long:"version" description:"显示版本号"`
		Install                bool   `short:"i" long:"install" description:"安装为系统服务"`
		Uninstall              bool   `long:"uninstall" description:"删除系统服务"`
		ShowConsole            bool   `long:"show-console" description:"Windows上显示控制台界面"`
		HideUIWhenBoot         bool   `long:"hide-ui" description:"启动时不弹出UI"`
		ServiceUser            string `long:"service-user" description:"用于启动服务的用户"`
		ServiceName            string `long:"service-name" description:"自定义服务名，默认为sealdice"`
		MultiInstanceOnWindows bool   `short:"m" long:"multi-instance" description:"允许在Windows上运行多个海豹"`
		Address                string `long:"address" description:"将UI的http服务地址改为此值，例: 0.0.0.0:3211"`
		DoUpdateWin            bool   `long:"do-update-win" description:"windows自动升级用，不要在任何情况下主动调用"`
		DoUpdateOthers         bool   `long:"do-update-others" description:"linux/mac自动升级用，不要在任何情况下主动调用"`
		Delay                  int64  `long:"delay"`
		JustForTest            bool   `long:"just-for-test"`
		DBCheck                bool   `long:"db-check" description:"检查数据库是否有问题"`
		ShowEnv                bool   `long:"show-env" description:"显示环境变量"`
		VacuumDB               bool   `long:"vacuum" description:"对数据库进行整理, 使其收缩到最小尺寸"`
		UpdateTest             bool   `long:"update-test" description:"更新测试"`
		LogLevel               int8   `long:"log-level" description:"设置日志等级" default:"0" choice:"-1" choice:"0" choice:"1" choice:"2" choice:"3" choice:"4" choice:"5"`
		ContainerMode          bool   `long:"container-mode" description:"容器模式，该模式下禁用内置客户端"`
	}

	if opts.Version {
		fmt.Println(dice.VERSION.String())
		return
	}
	if opts.DBCheck {
		model.DBCheck("data/default")
		return
	}
	if opts.VacuumDB {
		model.DBVacuum()
		return
	}
	if opts.ShowEnv {
		for i, e := range os.Environ() {
			println(i, e)
		}
		return
	}
	if opts.Delay != 0 {
		fmt.Println("延迟启动", opts.Delay, "秒")
		time.Sleep(time.Duration(opts.Delay) * time.Second)
	}

	_ = os.MkdirAll("./data", 0o755)
	MainLoggerInit("./data/main.log", true)

	diceLogger.SetEnableLevel(zapcore.Level(opts.LogLevel))

	// 提早初始化是为了读取ServiceName
	diceManager := &dice.DiceManager{}

	if opts.ContainerMode {
		logger.Info("当前为容器模式，内置适配器与更新功能已被禁用")
		diceManager.ContainerMode = true
	}

	diceManager.LoadDice()
	diceManager.IsReady = true

	if opts.Address != "" {
		fmt.Println("由参数输入了服务地址:", opts.Address)
		diceManager.ServeAddress = opts.Address
	}

	cwd, _ := os.Getwd()
	fmt.Printf("%s %s\n", dice.APPNAME, dice.VERSION.String())
	fmt.Println("工作路径: ", cwd)

	if strings.HasPrefix(cwd, os.TempDir()) {
		// C:\Users\XXX\AppData\Local\Temp
		// C:\Users\XXX\AppData\Local\Temp\BNZ.627d774316768935
		tempDirWarn()
		return
	}

	useBuiltinUI := true

	// 删除遗留的shm和wal文件
	if !model.DBCacheDelete() {
		logger.Error("数据库缓存文件删除失败")
		showMsgBox("数据库缓存文件删除失败", "为避免数据损坏，拒绝继续启动。请检查是否启动多份程序，或有其他程序正在使用数据库文件！")
		return
	}
	// v150升级
	if !migrate.V150Upgrade() {
		return
	}

	if !opts.ShowConsole || opts.MultiInstanceOnWindows {
		hideWindow()
	}

	go dice.TryGetBackendURL()

	cleanUp := cleanupCreate(diceManager)
	defer dice.CrashLog()
	defer cleanUp()

	// 初始化核心
	diceManager.TryCreateDefault()
	diceManager.InitDice()
	go func() {
		// 每5分钟做一次新版本检查
		for {
			go CheckVersion(diceManager)
			time.Sleep(5 * time.Minute)
		}
	}()
	go RebootRequestListen(diceManager)
	go UpdateRequestListen(diceManager)
	go UpdateCheckRequestListen(diceManager)

	// 强制清理机制
	go (func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-interrupt
		cleanUp()
		time.Sleep(3 * time.Second)
		os.Exit(0)
	})()

	if opts.Address != "" {
		fmt.Println("由参数输入了服务地址:", opts.Address)
	}

	for _, d := range diceManager.Dice {
		go diceServe(d)
	}

	// pprof
	// go func() {
	//	http.ListenAndServe("0.0.0.0:8899", nil)
	// }()

	go uiServe(diceManager, opts.HideUIWhenBoot, useBuiltinUI)
	/*go func() {
		// 创建一个新的 WebView 实例
		w := webview.New(1024, 768, false, true)
		defer w.Destroy()

		// 设置窗口标题和大小
		w.SetTitle("Tempest Dice")
		w.SetSize(1024, 768, webview.HintNone)

		// 设置要显示的 URL
		d := diceManager
		w.Navigate(fmt.Sprintf("%s%s", "http://localhost:", d.ServeAddress))
		// 运行 WebView
		w.Run()
	}()*/
	initEndTime := time.Now().UnixMicro()
	fmt.Printf("%s%d%s", "初始化完成，耗时: ", int((initEndTime-initStartTime)/1000), "毫秒\n")
	// OOM分析工具
	// err = nil
	// err = http.ListenAndServe(":9090", nil)
	// if err != nil {
	// 	fmt.Printf("ListenAndServe: %s", err)
	// }

	// darwin 的托盘菜单似乎需要在主线程启动才能工作，调整到这里
	trayInit(diceManager)
}

func removeUpdateFiles() {
	// 无论原因，只要走到这里全部删除
	_ = os.Remove("./auto_update_ok")
	_ = os.Remove("./auto_update.exe")
	_ = os.Remove("./auto_updat3.exe")
	_ = os.Remove("./auto_update_ok")
	_ = os.Remove("./auto_update")
	_ = os.Remove("./_delete_me.exe")
	_ = os.RemoveAll("./update")
}

func diceServe(d *dice.Dice) {
	defer dice.CrashLog()
	if len(d.ImSession.EndPoints) == 0 {
		d.Logger.Infof("未检测到任何帐号，请先到“帐号设置”进行添加")
	}

	d.UIEndpoint = new(dice.EndPointInfo)
	d.UIEndpoint.Enable = true
	d.UIEndpoint.Platform = "UI"
	d.UIEndpoint.ID = "1"
	d.UIEndpoint.State = 1
	d.UIEndpoint.UserID = "UI:1000"
	d.UIEndpoint.Adapter = &dice.PlatformAdapterHTTP{Session: d.ImSession, EndPoint: d.UIEndpoint}
	d.UIEndpoint.Session = d.ImSession

	dice.TextMapCompatibleCheckAll(d)

	for _, _conn := range d.ImSession.EndPoints {
		if _conn.Enable {
			go func(conn *dice.EndPointInfo) {
				defer dice.ErrorLogAndContinue(d)

				switch conn.Platform {
				case "QQ":
					if conn.EndPointInfoBase.ProtocolType == "onebot" {
						pa := conn.Adapter.(*dice.PlatformAdapterGocq)
						if pa.BuiltinMode == "lagrange" {
							dice.LagrangeServe(d, conn, dice.LagrangeLoginInfo{
								IsAsyncRun: true,
							})
							return
						} else {
							dice.GoCqhttpServe(d, conn, dice.GoCqhttpLoginInfo{
								Password:         pa.InPackGoCqhttpPassword,
								Protocol:         pa.InPackGoCqhttpProtocol,
								AppVersion:       pa.InPackGoCqhttpAppVersion,
								IsAsyncRun:       true,
								UseSignServer:    pa.UseSignServer,
								SignServerConfig: pa.SignServerConfig,
							})
						}
					}
				}
			}(_conn)
		} else {
			_conn.State = 0 // 重置状态
		}
	}
}

func uiServe(dm *dice.DiceManager, hideUI bool, useBuiltin bool) {
	logger.Info("即将启动webui")
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "token"},
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	mimePatch()
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline'; img-src 'self' data: blob: *; style-src  'self' 'unsafe-inline' *; frame-src 'self' *;",
		// XFrameOptions:         "ALLOW-FROM https://captcha.go-cqhttp.org/",
	}))
	// X-Content-Type-Options: nosniff

	groupStatic := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/" {
				responseWriter := c.Response()
				responseWriter.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				responseWriter.Header().Set("Pragma", "no-cache")
				responseWriter.Header().Set("Expires", "0")
			}
			return next(c)
		}
	}
	e.Use(groupStatic)
	if useBuiltin {
		frontend, _ := fs.Sub(static.Frontend, "frontend")
		e.StaticFS("/", frontend)
	} else {
		e.Static("/", "./frontend_overwrite")
	}

	api.Bind(e, dm)
	e.HideBanner = true // 关闭banner，原因是banner图案会改变终端光标位置

	httpServe(e, dm, hideUI)
}

//
// func checkCqHttpExists() bool {
//	if _, err := os.Stat("./go-cqhttp"); err == nil {
//		return true
//	}
//	return false
// }

func mimePatch() {
	builtinMimeTypesLower := map[string]string{
		".css":  "text/css; charset=utf-8",
		".gif":  "image/gif",
		".htm":  "text/html; charset=utf-8",
		".html": "text/html; charset=utf-8",
		".jpg":  "image/jpeg",
		".js":   "application/javascript",
		".wasm": "application/wasm",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".xml":  "text/xml; charset=utf-8",
	}

	for k, v := range builtinMimeTypesLower {
		_ = mime.AddExtensionType(k, v)
	}
}
