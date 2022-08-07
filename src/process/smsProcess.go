package processes

import (
	"common"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"net"
	"utils"
)

type SmsProcessor struct {
}

// SendMesToAllUsers 向所有在线的用户发送消息
func (sp *SmsProcessor) SendMesToAllUsers(mes *common.Message) (err error) {

	//将受到的消息反序列化
	var smsMes common.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		logger.Error("smsMes unmarshal err, err=%s", err.Error())
		return
	}
	//fmt.Println("smsMes=",smsMes)
	/*
		必须序列化整个smsMes结构体，因为要发送一个消息实例，而不是一个string
	*/
	//data, err := json.Marshal(smsMes.Content)
	//data, err := json.Marshal(smsMes)

	//组装服务器转发聊天消息实例
	smsRespMes := common.SmsRespMes{
		User:    smsMes.User,
		Content: smsMes.Content,
	}
	data, err := json.Marshal(smsRespMes)
	if err != nil {
		logger.Error("smsRespMes marshal err, err=%s", err.Error())
		return
	}
	mesResp := common.Message{
		Type: common.SmsRespMesType,
		Data: string(data),
	}
	data, err = json.Marshal(mesResp)
	if err != nil {
		logger.Error("mesResp marshal err, err=%s", err.Error())
		return
	}

	//拿到所有在线的用户和其连接
	for id, up := range Usermgr.OnlineUsers {
		if id == smsMes.UserId {
			continue
		}
		err = sp.SendMesToUser(data, up.Conn)
	}
	return
}

// SendMesToUser 向所有在线的用户发送消息
func (sp *SmsProcessor) SendMesToUser(data []byte, conn net.Conn) (err error) {
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
