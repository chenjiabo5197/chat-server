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
	UserId   string //后添加参数，用于表明该Conn是哪位用户的连接
	UserName string
}

// NotifyOthersOnlineUser 通知其他用户当前用户上线/下线的消息
func (up *UserProcess) NotifyOthersOnlineUser(userId string) (err error) {
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
func (up *UserProcess) ServerProcessLogin(mes *common.Message) (loginUser *common.User, err error, isSuccess bool) {
	//将mes反序列化成LoginMes结构体
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		logger.Error("loginMes unmarshal err, err=%s", err.Error())
		return
	}
	//定义一个LoginResMes结构体
	var loginRespMes common.LoginRespMes

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
		loginRespMes.User = *user
		//通知其他用户有新用户上线,2.0下线，转为用户查询再返回
		//np := NotifyProcessor{}
		//err = np.NotifyOthersOnlineUser(user, 0)
		//for k := range Usermgr.OnlineUsers {
		//	loginRespMes.UsersId = append(loginRespMes.UsersId, k)
		//}
		isSuccess = true
		loginUser = user
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

	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	logger.Info("send to client login resp data=%s", data)
	return
}

// ServerProcessRigister 专门用于处理注册的函数
func (up *UserProcess) ServerProcessRegister(mes *common.Message) (err error) {
	//将mes反序列化成LoginMes结构体
	var registerMes common.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)

	if err != nil {
		logger.Error("RegisterMes unmarshal err, err=%s", err.Error())
		return
	}
	logger.Info("RegisterMes=%s", utils.Struct2String(registerMes))
	//定义一个状态返回的结构体，表示执行的状态
	var registerRespMes common.StatusRespMes

	//拿UserDao对象的方法去redis验证
	err = model.MyUserDao.RegisterUser(registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerRespMes.RespCode = 500
			registerRespMes.Error = err.Error()
		} else {
			registerRespMes.RespCode = 505
			registerRespMes.Error = "服务器内部错误"
		}
	} else {
		registerRespMes.RespCode = 200
		logger.Info("register success")
	}
	//将rigisterResMes反序列化
	data, err := json.Marshal(registerRespMes)
	if err != nil {
		logger.Error("registerRespMes marshal err, err=%s", err.Error())
		return
	}
	err = up.SendRespStatus(string(data), common.RegisterRespMesType)
	return
}

// 用于将服务器对客户端的返回发送
func (up *UserProcess) SendRespStatus(sendData string, mesType common.MesType) (err error) {
	//定义一个Message结构体，用于发送给客户端
	var resMes common.Message
	resMes.Type = mesType
	resMes.Data = sendData

	//将Message结构体序列化
	data, err := json.Marshal(resMes)
	if err != nil {
		logger.Info("resMes marshal err, err=%v", err.Error())
		return
	}
	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		logger.Error("fail send to client status resp data=%s||err=%s", data, err.Error())
	}
	logger.Info("success send to client status resp data=%s", data)
	return
}
