package dice

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func RandString(n int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

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
