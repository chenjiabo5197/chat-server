package processes

import (
	"common"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"model"
	"net"
	"utils"
)

// UserProcess 用于处理用户的结构体(用户登录和注册)
type UserProcess struct {
	Conn     net.Conn
	UserId   int //后添加参数，用于表明该Conn是哪位用户的连接
	UserName string
}

// NotifyOthersOnlineUser 通知其他用户当前用户上线/下线的消息
func (up *UserProcess) NotifyOthersOnlineUser(userId int) (err error) {
	//实例一个返回消息对象
	var notifyMes common.NotifyUserStatusMes
	notifyMes.UserId = userId
	notifyMes.UserStatus = common.UserOnline

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
	//遍历onlineUsers 这个map，向其中的user通知上线信息
	for k, newUp := range Usermgr.OnlineUsers {
		if k == userId { //不向自己通知自己上线的消息
			continue
		}

		//开始通知
		err = newUp.NotifyOnlineUser(data, newUp.Conn)
	}
	return
}

// NotifyOnlineUser 通知其他用户当前用户上线/下线的消息  具体实现方法
func (up *UserProcess) NotifyOnlineUser(data []byte, conn net.Conn) (err error) {

	tf := utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		logger.Error("notifyOnlineUser err, err=%s", err.Error())
	}
	return
}

// ServerProcessLogin 专门用于处理登陆的函数
func (up *UserProcess) ServerProcessLogin(mes *common.Message) (userId int, err error) {
	//将mes反序列化成LoginMes结构体
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		logger.Error("loginMes unmarshal err, err=%s", err.Error())
		return
	}
	//定义一个LoginResMes结构体
	var loginRespMes common.LoginRespMes

	//判断输入的用户名和密码是否符合规定
	// if loginMes.UserId == 100 && loginMes.UserPwd == "abc" {
	// 	loginResMes.ResCode = 200
	// }else{
	// 	loginResMes.ResCode = 500
	// 	loginResMes.Error = "输入的用户id或密码不正确"
	// }

	//拿UserDao对象的方法去redis验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		//loginResMes.ResCode = 500
		//loginResMes.Error = "该用户不存在"
		if err == model.ERROR_USER_NOTEXISTS {
			loginRespMes.RespCode = 500
			loginRespMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginRespMes.RespCode = 403
			loginRespMes.Error = err.Error()
		} else {
			loginRespMes.RespCode = 505
			loginRespMes.Error = "服务器内部错误"
		}
	} else {
		//用户登录成功将用户加入到在线用户列表中
		up.UserName = user.UserName
		up.UserId = user.UserId
		Usermgr.AddOnlineUsers(up)
		loginRespMes.RespCode = 200
		loginRespMes.UserName = user.UserName
		loginRespMes.UsersId = user.UserId
		//通知其他用户有新用户上线,2.0下线，转为用户查询再返回
		//np := NotifyProcessor{}
		//err = np.NotifyOthersOnlineUser(user, 0)
		//for k := range Usermgr.OnlineUsers {
		//	loginRespMes.UsersId = append(loginRespMes.UsersId, k)
		//}
		logger.Info("%s登陆成功", utils.Struct2String(user))
	}

	//将LoginResMes反序列化
	data, err := json.Marshal(loginRespMes)
	if err != nil {
		logger.Error("LoginRespMes marshal err, err=%s", err.Error())
		return
	}

	//定义一个Message结构体，用于发送给客户端
	var resMes common.Message
	resMes.Type = common.LoginRespMesType
	resMes.Data = string(data)

	//将Message结构体序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		logger.Error("Message marshal err, err=%s", err.Error())
		return
	}

	logger.Info("send to client login resp data=%s", data)

	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return loginMes.UserId, err

}

// ServerProcessRigister 专门用于处理注册的函数
func (up *UserProcess) ServerProcessRigister(mes *common.Message) (err error) {
	//将mes反序列化成LoginMes结构体
	var rigisterMes common.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &rigisterMes)

	if err != nil {
		logger.Error("RigisterMes unmarshal err, err=%s", err.Error())
		return
	}
	logger.Info("RigisterMes=%s", utils.Struct2String(rigisterMes))
	//定义一个RigisterResMes结构体
	var rigisterRespMes common.RegisterRespMes

	//拿UserDao对象的方法去redis验证
	err = model.MyUserDao.RegisterUser(rigisterMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			rigisterRespMes.RespCode = 500
			rigisterRespMes.Error = err.Error()
		} else {
			rigisterRespMes.RespCode = 505
			rigisterRespMes.Error = "服务器内部错误"
		}
	} else {
		rigisterRespMes.RespCode = 200
		logger.Info("rigister success")
	}

	//将rigisterResMes反序列化
	data, err := json.Marshal(rigisterRespMes)
	if err != nil {
		logger.Error("rigisterRespMes marshal err, err=%s", err.Error())
		return
	}

	//定义一个Message结构体，用于发送给客户端
	var resMes common.Message
	resMes.Type = common.RegisterRespMesType
	resMes.Data = string(data)

	//将Message结构体序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		logger.Info("resMes marshal err, err=%v", err.Error())
		return
	}
	logger.Info("send to client rigister resp data=%s", data)

	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}
