package processes

import (
	"common"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"model"
	"net"
	"utils"
)

type SmsProcessor struct {
}

// SendMesToAllUsers 向所有在线的用户发送消息
func (sp *SmsProcessor) SendMesToAllUsers(mes *common.Message) (err error, statusStr string) {
	statusResp := common.StatusRespMes{}
	defer func() {
		if err != nil && statusResp.RespCode == 0 {
			statusResp.RespCode = 500
			statusResp.Error = err.Error()
		}
		statusStr = utils.Struct2String(statusResp)
	}()
	//将受到的消息反序列化
	var smsMes common.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	logger.Info("receive send to group message, data=%s", utils.Struct2String(smsMes))
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
		Type: common.RecvSmsMesType,
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

// SendMesToOne 向单个用户发送消息  返回状态，表示发送1对1消息的结果，失败或成功
func (sp *SmsProcessor) SendMesToOne(smsStr string) (err error, statusStr string) {
	statusResp := common.StatusRespMes{}
	defer func() {
		if err != nil && statusResp.RespCode == 0 {
			statusResp.RespCode = 500
			statusResp.Error = err.Error()
		}
		statusStr = utils.Struct2String(statusResp)
	}()
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
		Type: common.RecvSmsToOneMesType,
		Data: string(data),
	}
	data, err = json.Marshal(mesResp)
	if err != nil {
		logger.Error("mesResp marshal err, err=%s", err.Error())
		return
	}
	// 查找是否有这个用户
	targetUserId := utils.GetMd5Value(smsMes.SmsMesTarget)
	targetKey := model.GetRedisUserKey(targetUserId)
	_, err = model.MyUserDao.GetDataByKey(targetKey)
	if err != nil {  // 该用户不存在
		logger.Error("SendMesToOne||cannot find %s in redis", smsMes.SmsMesTarget)
		statusResp.Error = err.Error()
		statusResp.RespCode = 404
		return
	}
	// 在redis中找到该用户，用户存在，先判断是否在线，在线发送在线消息，不在线发送的消息储存在redis中，等上线后再发送过去
	isOnline := false
	for _, up := range Usermgr.OnlineUsers {
		if smsMes.SmsMesTarget == up.UserName {
			isOnline = true
			err = sp.SendMesToUser(data, up.Conn)
			logger.Info("success send to online one message, data=%s", string(data))
			break
		}
	}
	if !isOnline { // 该用户不在线，将消息存入redis
		err = model.MyUserDao.HSetDataByName(smsMes.UserName, mesResp)
		if err != nil {
			logger.Error("offline message hset to redis err, err=%v||mesResp=%s", err, utils.Struct2String(mesResp))
			return
		}
		logger.Info("offline message hset success store to redis, mesResp=%s", utils.Struct2String(mesResp))
	}
	logger.Info("success send to one message, data=%s||isOnline=%v", utils.Struct2String(mesResp), isOnline)
	return
}

// SendMesToUser 向在线用户发送消息
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
