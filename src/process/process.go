package processes

import (
	"common"
	logger "github.com/shengkehua/xlog4go"
	"io"
	"net"
	"utils"
)

// Processor 创建一个Process的结构体
type Processor struct {
	Conn     net.Conn
	CurId    string //保存此进程中和服务器连接的客户id，用于客户端下线后删除该客户端
	UserName string // 登陆的用户名
}

// ServerProcessMes 根据客户端传来的不同消息类型来调用不同的函数处理
func (p *Processor) ServerProcessMes(mes *common.Message) (err error) {
	//用于调试
	//fmt.Println("mes=",mes)

	switch mes.Type {
		case common.LoginMesType:  //处理登陆的函数
			up := UserProcess{
				Conn: p.Conn,
			}
			id, err := up.ServerProcessLogin(mes)
			if err == nil {
				//此时客户端已经登录成功
				logger.Info("%s login success", id)
			}
			return err
		case common.RegisterMesType:   //处理注册的函数
			up := UserProcess{
				Conn: p.Conn,
			}
			err = up.ServerProcessRegister(mes)
			return
		case common.SmsMesType:  //处理转发广播消息
			sp := SmsProcessor{}
			_, statusRespStr := sp.SendMesToAllUsers(mes)
			up := UserProcess{
				Conn: p.Conn,
			}
			// 将发送消息的结果返回给客户端
			err = up.SendRespStatus(statusRespStr, common.SmsRespMesType)
			return
		case common.QueryAllOnlineType:  //处理查询在线用户
			qo := QueryOnline{
				Conn: p.Conn,
			}
			qo.QueryAllOnlineUser(mes.Data)
			return
		case common.SmsToOneMesType:  //处理转发1对1聊天消息
			sp := SmsProcessor{}
			_, statusRespStr := sp.SendMesToOne(mes.Data)
			up := UserProcess{
				Conn: p.Conn,
			}
			// 将发送消息的结果返回给客户端
			err = up.SendRespStatus(statusRespStr, common.SmsRespMesType)
			return
		default:
			//错误
			logger.Error("unknown mes type，type=%v", mes.Type)
			return
	}
}

func (p *Processor) HandlerRecvMes() (err error) {
	//循环读取客户端发送的消息
	for {
		//封装函数，传一个conn连接，读取客户端的输入，返回mes和err
		//创建一个Transfer实例用于传输数据
		tf := utils.Transfer{
			Conn: p.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				logger.Info("client close connect")
				//客户端关闭连接，在在线用户列表中将其删除,再通知客户端在其本地维护的onlineUserMap中删除该用户
				Usermgr.DeleteOnlineUsers(p.CurId)
				np := NotifyProcessor{}
				user := &common.User{
					UserId: p.CurId,
				}
				err = np.NotifyOthersOnlineUser(user, 1)
				logger.Info("online user list=%s", utils.Struct2String(Usermgr.OnlineUsers))
				return err
			} else {
				logger.Error("read client list err, err=%s", err.Error())
			}
		}
		// fmt.Println("mes=",mes)
		err = p.ServerProcessMes(&mes)
		if err != nil {
			logger.Error("handler err, err=%s", err.Error())
			return err
		}
	}
}
