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
	logger.Info("receive send to group message, data=%s", utils.Struct2String(smsMes))
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
		User:        smsMes.User,
		Content:     smsMes.Content,
		SmsRespFrom: smsMes.UserName,
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
			continue // 不向自己发送
		}
		err = sp.SendMesToUser(data, up.Conn)
	}
	logger.Info("success send to group message, data=%s", utils.Struct2String(mesResp))
	return
}

// SendMesToOne 向单个用户发送消息
func (sp *SmsProcessor) SendMesToOne(smsStr string) (err error) {
	//将收到的消息反序列化
	var smsMes common.SmsMes
	err = json.Unmarshal([]byte(smsStr), &smsMes)
	logger.Info("receive send to one message, data=%s", utils.Struct2String(smsMes))
	if err != nil {
		logger.Error("smsMes unmarshal err, err=%s", err.Error())
		return
	}

	//组装服务器转发聊天消息实例
	smsRespMes := common.SmsRespMes{
		User:        smsMes.User,
		Content:     smsMes.Content,
		SmsRespFrom: smsMes.UserName,
	}
	data, err := json.Marshal(smsRespMes)
	if err != nil {
		logger.Error("smsRespMes marshal err, err=%s", err.Error())
		return
	}
	mesResp := common.Message{
		Type: common.SmsToOneRespMesType,
		Data: string(data),
	}
	data, err = json.Marshal(mesResp)
	if err != nil {
		logger.Error("mesResp marshal err, err=%s", err.Error())
		return
	}

	//遍历所有在线的用户和其连接，找到要发送的对象
	for _, up := range Usermgr.OnlineUsers {
		if smsMes.SmsMesTarget == up.UserName {
			err = sp.SendMesToUser(data, up.Conn)
			break
		}
	}
	logger.Info("success send to one message, data=%s", utils.Struct2String(mesResp))
	return
}

// SendMesToUser 向用户发送消息
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
