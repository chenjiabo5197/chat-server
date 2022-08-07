package processes

import (
	"common"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"net"
	"utils"
)

/*
	用于向其他客户端通知有客户端上线和下线的消息
*/
type NotifyProcessor struct {
}

// NotifyOthersOnlineUser 通知其他用户当前用户上线/下线的消息
func (np *NotifyProcessor) NotifyOthersOnlineUser(userId int, status int) (err error) {

	//实例化一个NotifyUserStatusMes对象，用于向客户端通知
	var notifyMes common.NotifyUserStatusMes
	notifyMes.UserId = userId
	switch status {
	case 0:
		notifyMes.UserStatus = common.UserOnline
	case 1:
		notifyMes.UserStatus = common.UserOffline
	default:
		logger.Error("unknown type")
		return
	}

	data, err := json.Marshal(notifyMes)
	if err != nil {
		logger.Error("notifyMes marshal err, err=%s", err.Error())
		return
	}

	//再定义一个返回消息的实例对象
	var mes common.Message
	mes.Type = common.NotifyUserStatusMesType
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		logger.Error("mes marshal err, err=%s", err.Error())
		return
	}

	if status == 0 { //通知上线信息
		//先遍历拿到在线用户列表
		for id, newUp := range Usermgr.OnlineUsers {
			if id == userId { //跳过自己，不向自己通知自己上线的消息
				continue
			}
			err = np.NotifyOnlineUser(data, newUp.Conn)
			return
		}
	} else { //通知离线消息
		//先遍历拿到在线用户列表
		for _, newUp := range Usermgr.OnlineUsers {
			//此处不用判断是否给自己传递，因为调用函数前已经将自己从在线列表删除
			err = np.NotifyOnlineUser(data, newUp.Conn)
			return
		}
	}

	return
}

/*
	通知其他用户当前用户上线/下线的消息  具体实现方法
*/
func (np *NotifyProcessor) NotifyOnlineUser(data []byte, conn net.Conn) (err error) {

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		logger.Error("service transfer mes err, err=%s", err.Error())
		return
	}
	return

}
