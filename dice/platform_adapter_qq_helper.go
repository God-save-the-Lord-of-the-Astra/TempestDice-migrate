package dice

import (
	"time"
)

func ServeQQ(d *Dice, ep *EndPointInfo) {
	defer CrashLog()
	if ep.Platform != "QQ" {
		return
	}

	switch ep.ProtocolType {
	case "onebot":
		fallthrough
	default: // onebot 作为默认情况
		conn := ep.Adapter.(*PlatformAdapterGocq)
		serverGocq(d, ep, conn)
	}
}

func serverGocq(d *Dice, ep *EndPointInfo, conn *PlatformAdapterGocq) {
	if conn.diceServing {
		return
	}
	conn.diceServing = true

	ep.Enable = true
	ep.State = 2 // 连接中
	d.LastUpdatedTime = time.Now().Unix()
	d.Save(false)

	checkQuit := func() bool {
		if conn.GoCqhttpState == StateCodeInLoginDeviceLock {
			d.Logger.Infof("检测到设备锁流程，暂时不再连接")
			ep.State = 0
			d.LastUpdatedTime = time.Now().Unix()
			d.Save(false)
			return true
		}
		if !conn.diceServing {
			// 退出连接
			d.Logger.Infof("检测到连接关闭，不再进行此onebot服务的重连: <%s>(%s)", ep.Nickname, ep.UserID)
			return true
		}
		if conn.GoCqhttpState == StateCodeLoginFailed {
			d.Logger.Infof("检测到登录失败，不再进行此onebot服务的重连: <%s>(%s)", ep.Nickname, ep.UserID)
			return true
		}
		if !ep.Enable {
			d.Logger.Infof("检测到账号被禁用，不再进行此onebot服务的重连: <%s>(%s)", ep.Nickname, ep.UserID)
			return true
		}
		return false
	}

	conn.reconnectTimes = 0
	for {
		if checkQuit() {
			break
		}

		// 骰子开始连接
		d.Logger.Infof("开始连接 onebot 服务，帐号 <%s>(%s)，重试计数[%d/%d]", ep.Nickname, ep.UserID, conn.reconnectTimes, 5)
		ret := ep.Adapter.Serve()

		if ret == 0 {
			break
		}

		if checkQuit() {
			break
		}

		if conn.GoCqhttpState == StateCodeInLogin || conn.GoCqhttpState == StateCodeInLoginQrCode {
			time.Sleep(15 * time.Second)
			continue
		}

		conn.reconnectTimes++
		if conn.reconnectTimes > 5 {
			d.Logger.Infof("onebot 连接重试次数过多，先行中断: <%s>(%s)", ep.Nickname, ep.UserID)
			ep.State = 0
			conn.GoCqhttpState = StateCodeLoginFailed
			break
		}

		time.Sleep(15 * time.Second)
	}

	conn.diceServing = false
}
