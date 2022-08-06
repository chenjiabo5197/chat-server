package processes

import (
	"common"
	"fmt"
	"io"
	"net"
	"utils"
)

//创建一个Process的结构体
type Processor struct {
	Conn  net.Conn
	CurId int //保存此进程中和服务器连接的客户id，用于客户端下线后删除该客户端
}

/*
	根据客户端传来的不同消息类型来调用不同的函数处理
*/
func (p *Processor) ServerProcessMes(mes *common.Message) (err error) {
	//用于调试
	//fmt.Println("mes=",mes)

	switch mes.Type {
	case common.LoginMesType:
		//处理登陆的函数
		up := UserProcess{
			Conn: p.Conn,
		}
		id, err := up.ServerProcessLogin(mes)
		if err == nil {
			//此时客户端已经登录成功
			p.CurId = id
			//fmt.Println("p.CurId=",p.CurId)
		}
		return err
	case common.RegisterMesType:
		//处理注册的函数
		up := UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRigister(mes)
		return
	case common.SmsMesType:
		//处理转发聊天消息
		sp := SmsProcessor{}
		err = sp.SendMesToAllUsers(mes)
		return
	default:
		//错误
		fmt.Println("错误的消息类型")
		return
	}
}

func (p *Processor) Process2() (err error) {
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
				fmt.Println("客户端关闭了连接")
				//客户端关闭连接，在在线用户列表中将其删除,再通知客户端在其本地维护的onlineUserMap中删除该用户
				Usermgr.DeleteOnlineUsers(p.CurId)
				np := NotifyProcessor{}
				err = np.NotifyOthersOnlineUser(p.CurId, 1)
				fmt.Println("在线用户OnlineUsers=", Usermgr.OnlineUsers)
				return err
			} else {
				fmt.Println("读取客户端消失失败,err=", err)
			}
		}
		// fmt.Println("mes=",mes)
		err = p.ServerProcessMes(&mes)
		if err != nil {
			fmt.Println("err=", err)
			return err
		}
	}
}
