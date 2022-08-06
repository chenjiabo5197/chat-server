package processes

import (
	"common"
	"encoding/json"
	"fmt"
	"model"
	"net"
	"utils"
)

//用于处理用户的结构体(用户登录和注册)
type UserProcess struct {
	Conn   net.Conn
	UserId int //后添加参数，用于表明该Conn是哪位用户的连接
}

/*
	通知其他用户当前用户上线/下线的消息
*/
func (up *UserProcess) NotifyOthersOnlineUser(userId int) (err error) {
	//实例一个返回消息对象
	var notifyMes common.NotifyUserStatusMes
	notifyMes.UserId = userId
	notifyMes.UserStatus = common.UserOnline

	data, err := json.Marshal(notifyMes)
	if err != nil {
		fmt.Println("notifyMes序列化失败,err=", err)
		return
	}

	//再定义一个返回消息的实例对象
	var mes common.Message
	mes.Type = common.NotifyUserStatusMesType
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes反序列化失败,err=", err)
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

/*
	通知其他用户当前用户上线/下线的消息  具体实现方法
*/
func (up *UserProcess) NotifyOnlineUser(data []byte, conn net.Conn) (err error) {

	tf := utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("err=", err)
	}
	return
}

/*
	专门用于处理登陆的函数
*/
func (up *UserProcess) ServerProcessLogin(mes *common.Message) (userId int, err error) {
	//将mes反序列化成LoginMes结构体
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("将LoginMes反序列化失败,err=", err)
		return
	}
	//定义一个LoginResMes结构体
	var loginResMes common.LoginResMes

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
			loginResMes.ResCode = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.ResCode = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.ResCode = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		/*
			用户登录成功先将用户加入到在线用户列表中,然后返回在线用户列表
		*/
		loginResMes.ResCode = 200
		Usermgr.OnlineUsers[loginMes.UserId] = up
		//通知其他用户有新用户上线
		np := NotifyProcessor{}
		err = np.NotifyOthersOnlineUser(loginMes.UserId, 0)
		for k := range Usermgr.OnlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, k)
		}
		fmt.Println(user, "登陆成功")
	}

	//将LoginResMes反序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("loginResMes结构体反序列化失败")
		return
	}

	//定义一个Message结构体，用于发送给客户端
	var resMes common.Message
	resMes.Type = common.LoginResMesType
	resMes.Data = string(data)

	//将Message结构体序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("Message序列化出错,err=", err)
		return
	}

	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return loginMes.UserId, err

}

/*
	专门用于处理注册的函数
*/
func (up *UserProcess) ServerProcessRigister(mes *common.Message) (err error) {
	//将mes反序列化成LoginMes结构体
	var rigisterMes common.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &rigisterMes)

	if err != nil {
		fmt.Println("将LoginMes反序列化失败,err=", err)
		return
	}
	fmt.Println("rigisterMes=", rigisterMes)
	//定义一个RigisterResMes结构体
	var rigisterResMes common.RegisterResMes

	//拿UserDao对象的方法去redis验证
	err = model.MyUserDao.RegisterUser(rigisterMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			rigisterResMes.ResCode = 500
			rigisterResMes.Error = err.Error()
		} else {
			rigisterResMes.ResCode = 505
			rigisterResMes.Error = "服务器内部错误"
		}
	} else {
		rigisterResMes.ResCode = 200
		fmt.Println("注册成功")
	}

	//将LoginResMes反序列化
	data, err := json.Marshal(rigisterResMes)
	if err != nil {
		fmt.Println("rigisterResMes结构体反序列化失败")
		return
	}

	//定义一个Message结构体，用于发送给客户端
	var resMes common.Message
	resMes.Type = common.RegisterResMesType
	resMes.Data = string(data)

	//将Message结构体序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("Message序列化出错,err=", err)
		return
	}

	//将Message序列化的结果发送给客户端
	//mvc模式，所以必须先创建一个Transfer实例，然后读取
	tf := utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}
