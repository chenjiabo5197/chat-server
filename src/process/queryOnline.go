package processes

import (
	"common"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"net"
	"utils"
)

type QueryOnline struct {
	Conn net.Conn
}

func (qo *QueryOnline) QueryAllOnlineUser(queryStr string) {
	//将请求的用户信息解析出来
	user := common.User{}
	err := json.Unmarshal([]byte(queryStr), &user)
	if err != nil {
		logger.Error("unmarshal QueryOnlineStr err, err=%s", err.Error())
		return
	}
	onlineUserName := make([]string, 0)
	logger.Debug("online User data=%s", utils.Struct2String(Usermgr.OnlineUsers))
	//拿到所有在线的用户和其连接
	for id, up := range Usermgr.OnlineUsers {
		if id == user.UserId {
			continue
		}
		onlineUserName = append(onlineUserName, up.UserName)
	}
	mes := common.Message{}
	mes.Type = common.AllOnlineRespType
	mes.Data = utils.Struct2String(onlineUserName)
	mesByte, _ := json.Marshal(mes)
	tf := utils.Transfer{
		Conn: qo.Conn,
	}
	err = tf.WritePkg(mesByte)
	if err != nil {
		logger.Error("send to client err, err=%s", err.Error())
		return
	}
	logger.Info("QueryAllOnlineUser|| send Online User data=%s", utils.Struct2String(mes))
}
