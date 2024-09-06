package dice

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type deviceFile struct {
	Display      string         `json:"display"`
	Product      string         `json:"product"`
	Device       string         `json:"device"`
	Board        string         `json:"board"`
	Model        string         `json:"model"`
	FingerPrint  string         `json:"finger_print"`
	BootID       string         `json:"boot_id"`
	ProcVersion  string         `json:"proc_version"`
	Protocol     int            `json:"protocol"` // 0: iPad 1: Android 2: AndroidWatch  // 3 macOS 4 企点
	IMEI         string         `json:"imei"`
	Brand        string         `json:"brand"`
	Bootloader   string         `json:"bootloader"`
	BaseBand     string         `json:"base_band"`
	SimInfo      string         `json:"sim_info"`
	OSType       string         `json:"os_type"`
	MacAddress   string         `json:"mac_address"`
	IPAddress    []int32        `json:"ip_address"`
	WifiBSSID    string         `json:"wifi_bssid"`
	WifiSSID     string         `json:"wifi_ssid"`
	ImsiMd5      string         `json:"imsi_md5"`
	AndroidID    string         `json:"android_id"`
	APN          string         `json:"apn"`
	VendorName   string         `json:"vendor_name"`
	VendorOSName string         `json:"vendor_os_name"`
	Version      *osVersionFile `json:"version"`
}

type osVersionFile struct {
	Incremental string `json:"incremental"`
	Release     string `json:"release"`
	Codename    string `json:"codename"`
	Sdk         uint32 `json:"sdk"`
}

func randomMacAddress() string {
	buf := make([]byte, 6)
	_, err := crand.Read(buf)
	if err != nil {
		return "00:16:ea:ae:3c:40"
	}
	// Set the local bit
	buf[0] |= 2
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func RandString(n int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// model	设备
// "iPhone11,2"	iPhone XS
// "iPhone11,8"	iPhone XR
// "iPhone12,1"	iPhone 11
// "iPhone13,2"	iPhone 12
// "iPad8,1"	iPad Pro
// "iPad11,2"	iPad mini
// "iPad13,2"	iPad Air 4
// "Apple Watch"	Apple Watch

var defaultConfig = `
# go-cqhttp 默认配置文件

account: # 账号相关
  uin: {QQ帐号} # QQ账号
  password: {QQ密码} # 密码为空时使用扫码登录
  encrypt: false  # 是否开启密码加密
  status: 0      # 在线状态 请参考 https://docs.go-cqhttp.org/guide/config.html#在线状态
  relogin: # 重连设置
    delay: 3   # 首次重连延迟, 单位秒
    interval: 3   # 重连间隔
    max-times: 0  # 最大重连次数, 0为无限制

  # 是否使用服务器下发的新地址进行重连
  # 注意, 此设置可能导致在海外服务器上连接情况更差
  use-sso-address: true
  # 是否允许发送临时会话消息
  allow-temp-session: false
{旧版签名服务相关配置信息}
{新版签名服务相关配置信息}

heartbeat:
  # 心跳频率, 单位秒
  # -1 为关闭心跳
  interval: 5

message:
  # 上报数据类型
  # 可选: string,array
  post-format: string
  # 是否忽略无效的CQ码, 如果为假将原样发送
  ignore-invalid-cqcode: false
  # 是否强制分片发送消息
  # 分片发送将会带来更快的速度
  # 但是兼容性会有些问题
  force-fragment: false
  # 是否将url分片发送
  fix-url: false
  # 下载图片等请求网络代理
  proxy-rewrite: ''
  # 是否上报自身消息
  report-self-message: false
  # 移除服务端的Reply附带的At
  remove-reply-at: false
  # 为Reply附加更多信息
  extra-reply-data: false
  # 跳过 Mime 扫描, 忽略错误数据
  skip-mime-scan: false

output:
  # 日志等级 trace,debug,info,warn,error
  log-level: warn
  # 日志时效 单位天. 超过这个时间之前的日志将会被自动删除. 设置为 0 表示永久保留.
  log-aging: 15
  # 是否在每次启动时强制创建全新的文件储存日志. 为 false 的情况下将会在上次启动时创建的日志文件续写
  log-force-new: true
  # 是否启用日志颜色
  log-colorful: true
  # 是否启用 DEBUG
  debug: false # 开启调试模式

# 默认中间件锚点
default-middlewares: &default
  # 访问密钥, 强烈推荐在公网的服务器设置
  access-token: ''
  # 事件过滤器文件目录
  filter: ''
  # API限速设置
  # 该设置为全局生效
  # 原 cqhttp 虽然启用了 rate_limit 后缀, 但是基本没插件适配
  # 目前该限速设置为令牌桶算法, 请参考:
  # https://baike.baidu.com/item/%E4%BB%A4%E7%89%8C%E6%A1%B6%E7%AE%97%E6%B3%95/6597000?fr=aladdin
  rate-limit:
    enabled: false # 是否启用限速
    frequency: 1  # 令牌回复频率, 单位秒
    bucket: 1     # 令牌桶大小

database: # 数据库相关设置
  leveldb:
    # 是否启用内置leveldb数据库
    # 启用将会增加10-20MB的内存占用和一定的磁盘空间
    # 关闭将无法使用 撤回 回复 get_msg 等上下文相关功能
    enable: true

  # 媒体文件缓存， 删除此项则使用缓存文件(旧版行为)
  cache:
    image: data/image.db
    video: data/video.db

# 连接服务列表
servers:
  # 添加方式，同一连接方式可添加多个，具体配置说明请查看文档
  #- http: # http 通信
  #- ws:   # 正向 Websocket
  #- ws-reverse: # 反向 Websocket
  #- pprof: #性能分析服务器
  # 正向WS设置
  - ws:
      # 正向WS服务器监听地址
      host: 127.0.0.1
      # 正向WS服务器监听端口
      port: {WS端口}
      # rc3
      address: 127.0.0.1:{WS端口}
      middlewares:
        <<: *default # 引用默认中间件
`

func NewGoCqhttpConnectInfoItem(account string) *EndPointInfo {
	conn := new(EndPointInfo)
	conn.ID = uuid.New().String()
	conn.Platform = "QQ"
	conn.ProtocolType = "onebot"
	conn.Enable = false
	conn.RelWorkDir = "extra/go-cqhttp-qq" + account

	conn.Adapter = &PlatformAdapterGocq{
		EndPoint:        conn,
		UseInPackClient: true,
		BuiltinMode:     "gocq",
	}
	return conn
}

func BuiltinQQServeProcessKillBase(dice *Dice, conn *EndPointInfo, isSync bool) {
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				dice.Logger.Error("内置 QQ 客户端清理报错: ", r)
				// go-cqhttp/lagrange 进程退出: exit status 1
			}
		}()

		pa, ok := conn.Adapter.(*PlatformAdapterGocq)
		if !ok {
			return
		}
		if !pa.UseInPackClient {
			return
		}

		// 重置状态
		conn.State = 0
		pa.GoCqhttpState = 0
		pa.GoCqhttpQrcodeData = nil

		if pa.BuiltinMode == "lagrange" {
			workDir := lagrangeGetWorkDir(dice, conn)
			qrcodeFile := filepath.Join(workDir, fmt.Sprintf("qr-%s.png", conn.UserID[3:]))
			if _, err := os.Stat(qrcodeFile); err == nil {
				// 如果已经存在二维码文件，将其删除
				_ = os.Remove(qrcodeFile)
				dice.Logger.Info("onebot: 删除已存在的二维码文件")
			}

			// 注意这个会panic，因此recover捕获了
			if pa.GoCqhttpProcess != nil {
				p := pa.GoCqhttpProcess
				pa.GoCqhttpProcess = nil
				// sigintwindows.SendCtrlBreak(p.Cmds[0].Process.Pid)
				_ = p.Stop()
				_ = p.Wait() // 等待进程退出，因为Stop内部是Kill，这是不等待的
			}
		} else {
			pa.GoCqhttpLoginDeviceLockURL = ""

			workDir := gocqGetWorkDir(dice, conn)
			qrcodeFile := filepath.Join(workDir, "qrcode.png")
			if _, err := os.Stat(qrcodeFile); err == nil {
				// 如果已经存在二维码文件，将其删除
				_ = os.Remove(qrcodeFile)
				dice.Logger.Info("onebot: 删除已存在的二维码文件")
			}

			// 注意这个会panic，因此recover捕获了
			if pa.GoCqhttpProcess != nil {
				p := pa.GoCqhttpProcess
				pa.GoCqhttpProcess = nil
				// sigintwindows.SendCtrlBreak(p.Cmds[0].Process.Pid)
				_ = p.Stop()
				_ = p.Wait() // 等待进程退出，因为Stop内部是Kill，这是不等待的
			}
		}
	}
	if isSync {
		f()
	} else {
		go f()
	}
}

func BuiltinQQServeProcessKill(dice *Dice, conn *EndPointInfo) {
	BuiltinQQServeProcessKillBase(dice, conn, false)
}

func gocqGetWorkDir(dice *Dice, conn *EndPointInfo) string {
	workDir := filepath.Join(dice.BaseConfig.DataDir, conn.RelWorkDir)
	return workDir
}

type GoCqhttpLoginInfo struct {
	UIN              int64
	Password         string
	Protocol         int
	AppVersion       string
	IsAsyncRun       bool
	UseSignServer    bool
	SignServerConfig *SignServerConfig
}

type SignServerConfig struct {
	SignServers          []*SignServer `yaml:"signServers" json:"signServers"`
	RuleChangeSignServer int           `yaml:"ruleChangeSignServer" json:"ruleChangeSignServer"`
	MaxCheckCount        int           `yaml:"maxCheckCount" json:"maxCheckCount"`
	SignServerTimeout    int           `yaml:"signServerTimeout" json:"signServerTimeout"`
	AutoRegister         bool          `yaml:"autoRegister" json:"autoRegister"`
	AutoRefreshToken     bool          `yaml:"autoRefreshToken" json:"autoRefreshToken"`
	RefreshInterval      int           `yaml:"refreshInterval" json:"refreshInterval"`
}

type SignServer struct {
	URL           string `yaml:"url" json:"url"`
	Key           string `yaml:"key" json:"key"`
	Authorization string `yaml:"authorization" json:"authorization"`
}

func GoCqhttpServe(dice *Dice, conn *EndPointInfo, loginInfo GoCqhttpLoginInfo) {
	pa := conn.Adapter.(*PlatformAdapterGocq)
	// if pa.GoCqHttpState != StateCodeInit {
	//	return
	//}
	pa.GoCqhttpState = StateCodeLoginSuccessed
	pa.GoCqhttpLoginSucceeded = true
	dice.Save(false)
	go ServeQQ(dice, conn)
}
