package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"vs_code/project0007/common"
	"vs_code/project0007/service/utils"
)

type SmsProcessor struct {

}

/*
	向所有在线的用户发送消息
 */
func (sp *SmsProcessor) SendMesToAllUsers (mes *common.Message) (err error) {

	//将受到的消息反序列化
	var smsMes common.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("smsMes反序列化失败,err=",err)
		return
	}
	//fmt.Println("smsMes=",smsMes)
	/*
		必须序列化整个smsMes结构体，因为要发送一个消息实例，而不是一个string
	 */
	//data, err := json.Marshal(smsMes.Content)
	//data, err := json.Marshal(smsMes)

	//组装服务器转发聊天消息实例
	smsResMes := common.SmsResMes{
		User:    smsMes.User,
		Content: smsMes.Content,
	}
	data, err := json.Marshal(smsResMes)
	if err != nil {
		fmt.Println("smsResMes序列化失败")
		return
	}
	mesRes := common.Message{
		Type: common.SmsResMesType,
		Data: string(data),
	}
	data, err = json.Marshal(mesRes)
	if err != nil {
		fmt.Println("mesRes序列化失败")
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


/*
	向所有在线的用户发送消息
*/
func (sp *SmsProcessor) SendMesToUser (data []byte, conn net.Conn) (err error) {
	tf := &utils.Transfer{
		Conn:conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("服务器转发消息失败")
		return
	}
	return

}
